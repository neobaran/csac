package lets

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"strings"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
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

	// 上传证书
	certId, err := sslClient.UploadCertificateData(string(resource.Certificate), string(resource.PrivateKey))
	if err != nil {
		return err
	}

	// 获取证书可用 CDN 域名
	domains, err := cdnClient.GetCDNDomains(certId)
	if err != nil {
		return err
	}

	// 更新 CDN 证书
	for _, item := range domains {
		_ = cdnClient.UpdateCDNConfig(item, certId)
	}

	return nil
}
