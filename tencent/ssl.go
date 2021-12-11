package tencent

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

type sslClient struct {
	*ssl.Client
}

// NewSSLClient 初始化SSL客户端
func (helper *TencentCloudHelper) NewSSLClient() *sslClient {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ssl.tencentcloudapi.com"
	client, _ := ssl.NewClient(helper.credential, "", cpf)

	return &sslClient{
		Client: client,
	}
}

// UploadCertificateData 上传证书信息
func (client *sslClient) UploadCertificateData(publicKey string, privateKey string) (*string, error) {
	request := ssl.NewUploadCertificateRequest()

	request.CertificatePublicKey = common.StringPtr(publicKey)
	request.CertificatePrivateKey = common.StringPtr(privateKey)

	response, err := client.UploadCertificate(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return response.Response.CertificateId, nil
}
