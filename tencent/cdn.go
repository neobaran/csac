package tencent

import (
	"fmt"

	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
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

func (client *cdnClient) GetCDNDomains(certId *string) ([]*string, error) {
	request := cdn.NewDescribeCertDomainsRequest()
	request.CertId = certId

	response, err := client.DescribeCertDomains(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s", response.ToJsonString())

	return response.Response.Domains, nil
}

func (client *cdnClient) UpdateCDNConfig(domain *string, certId *string) error {
	request := cdn.NewUpdateDomainConfigRequest()
	request.Domain = domain
	request.Https = &cdn.Https{
		CertInfo: &cdn.ServerCert{
			CertId: certId,
		},
	}

	_, err := client.UpdateDomainConfig(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return err
	}
	if err != nil {
		return err
	}

	return nil
}
