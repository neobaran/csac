package tencent

import (
	"github.com/neobaran/csac/config"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

type dnsPodClient struct {
	*dnspod.Client
}

// NewDNSPodClient 初始化dnspod客户端
func (helper *TencentCloudHelper) NewDNSPodClient() *dnsPodClient {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	client, _ := dnspod.NewClient(helper.credential, "", cpf)

	return &dnsPodClient{
		Client: client,
	}
}

// CreateRecordData 创建dns记录
func (client *dnsPodClient) CreateRecordData(domainData *Domain, value string, TTL uint64) error {
	request := dnspod.NewCreateRecordRequest()
	request.Domain = common.StringPtr(domainData.domain)
	request.SubDomain = common.StringPtr(domainData.subDomain)
	request.RecordType = common.StringPtr("TXT")
	request.RecordLine = common.StringPtr("默认")
	request.Value = common.StringPtr(value)
	request.TTL = common.Uint64Ptr(TTL)

	_, err := client.CreateRecord(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return err
	}
	if err != nil {
		return err
	}

	return nil
}

func (client *dnsPodClient) ListRecordData(domainData *Domain) ([]*dnspod.RecordListItem, error) {
	request := dnspod.NewDescribeRecordListRequest()
	request.Domain = common.StringPtr(domainData.domain)
	request.Subdomain = common.StringPtr(domainData.subDomain)
	request.RecordType = common.StringPtr("TXT")

	response, err := client.DescribeRecordList(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return response.Response.RecordList, nil
}

func (client *dnsPodClient) DeleteRecordData(domainData *Domain, item *dnspod.RecordListItem) error {
	request := dnspod.NewDeleteRecordRequest()
	request.Domain = common.StringPtr(domainData.domain)
	request.RecordId = item.RecordId
	_, err := client.DeleteRecord(request)

	return err
}

func (client *dnsPodClient) GetDNSProviderForLego(config config.Config) *DNSProvider {
	return NewDNSProvider(client, &config)
}
