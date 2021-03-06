package lets

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/log"
	"github.com/go-acme/lego/v4/registration"
	"github.com/neobaran/csac/config"
	"github.com/neobaran/csac/tencent"
)

var KeyTypes = map[string]certcrypto.KeyType{
	"EC256":   certcrypto.EC256,
	"EC384":   certcrypto.EC384,
	"RSA2048": certcrypto.RSA2048,
	"RSA4096": certcrypto.RSA4096,
	"RSA8192": certcrypto.RSA8192,
}

var wg sync.WaitGroup

type csac struct {
	cloudHelper *tencent.TencentCloudHelper
	lego        *lego.Client
}

func NewCSACHelper(config *config.Config, cloudHelper *tencent.TencentCloudHelper, debug bool) (*csac, error) {

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	user := User{
		Email: config.Email,
		key:   privateKey,
	}

	legoConfig := lego.NewConfig(&user)
	legoConfig.Certificate.KeyType = KeyTypes[strings.ToUpper(config.KeyType)]
	if debug {
		legoConfig.CADirURL = lego.LEDirectoryStaging
	}
	client, err := lego.NewClient(legoConfig)
	if err != nil {
		return nil, err
	}

	//init dns provider
	err = client.Challenge.SetDNS01Provider(cloudHelper.NewDNSPodClient().GetDNSProviderForLego(*config))
	if err != nil {
		return nil, err
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, err
	}
	user.Registration = reg

	return &csac{
		lego:        client,
		cloudHelper: cloudHelper,
	}, nil
}

func (t *csac) CreateSSL(domains []string) (*certificate.Resource, error) {
	request := certificate.ObtainRequest{
		Domains: domains,
		Bundle:  true,
	}
	certificates, err := t.lego.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}

	return certificates, nil
}

func (t *csac) UploadToCloud(resource *certificate.Resource) error {
	sslClient := t.cloudHelper.NewSSLClient()
	cdnClient := t.cloudHelper.NewCDNClient()

	// ????????????
	certId, err := sslClient.UploadCertificateData(string(resource.Certificate), string(resource.PrivateKey))
	if err != nil {
		return err
	}

	// ?????????????????? CDN ??????
	log.Infof("Start get CDN domain")
	domains, err := cdnClient.GetCDNDomains(certId, tencent.CDN)
	if err != nil {
		return err
	}
	time.Sleep(time.Second)
	// ?????????????????? ECDN ??????
	log.Infof("Start get ECDN domains")
	ecdnDomains, err := cdnClient.GetCDNDomains(certId, tencent.ECDN)
	if err != nil {
		return err
	}

	domains = append(domains, ecdnDomains...)
	log.Infof("[%s] Find %s domains for certificate %s", *certId, fmt.Sprint(len(domains)))

	// ?????? CDN ??????
	ch := make(chan struct{}, 10)
	for _, item := range domains {
		ch <- struct{}{}
		wg.Add(1)
		go func(domain *string) {
			defer wg.Done()
			log.Infof("[%s] Start to update CDN domain %s", *certId, *domain)
			_ = cdnClient.UpdateCDNConfig(domain, certId)
			time.Sleep(time.Second)
			<-ch
		}(item)
	}
	wg.Wait()

	return nil
}
