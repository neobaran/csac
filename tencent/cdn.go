package tencent

import (
	"github.com/go-acme/lego/v4/log"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const (
	CDN  = "cdn"
	ECDN = "ecdn"
)

type cdnClient struct {
	*cdn.Client
}

// NewCDNClient 初始化CDN客户端
func (helper *TencentCloudHelper) NewCDNClient() *cdnClient {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cdn.tencentcloudapi.com"
	client, _ := cdn.NewClient(helper.credential, "", cpf)

	return &cdnClient{
		Client: client,
	}
}

func (client *cdnClient) GetCDNDomains(certId *string, Product string) ([]*string, error) {
	request := cdn.NewDescribeCertDomainsRequest()
	request.CertId = certId
	request.Product = &Product

	response, err := client.DescribeCertDomains(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return response.Response.Domains, nil
}

func (client *cdnClient) UpdateCDNConfig(domain *string, certId *string) error {
	request := cdn.NewUpdateDomainConfigRequest()
	request.Domain = domain
	request.Https = &cdn.Https{
		Switch: common.StringPtr("on"),
		CertInfo: &cdn.ServerCert{
			CertId: certId,
		},
	}

	response, err := client.UpdateDomainConfig(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return err
	}
	if err != nil {
		return err
	}

	log.Infof("Update CDN domain [%s], response: %s", *domain, response.ToJsonString())
	return nil
}
