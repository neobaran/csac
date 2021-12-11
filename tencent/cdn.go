package tencent

import (
	"fmt"

	"github.com/go-acme/lego/v4/log"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
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

// GetCDNDomains 获取CDN域名
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

	return response.Response.Domains, nil
}

// UpdateCDNConfig 更新CDN配置
func (client *cdnClient) UpdateCDNConfig(domain *string, certId *string) error {

	searchRequest := cdn.NewDescribeDomainsConfigRequest()
	searchRequest.Limit = common.Int64Ptr(1)
	searchRequest.Filters = []*cdn.DomainFilter{
		{
			Name: common.StringPtr("domain"),
			Value: []*string{
				domain,
			},
			Fuzzy: common.BoolPtr(false),
		},
	}

	result, err := client.DescribeDomainsConfig(searchRequest)
	if err != nil {
		return err
	}

	if len(result.Response.Domains) > 0 {
		domainResult := result.Response.Domains[0]
		request := cdn.NewUpdateDomainConfigRequest()
		request.Domain = domainResult.Domain
		request.ProjectId = domainResult.ProjectId
		request.Origin = domainResult.Origin
		request.IpFilter = domainResult.IpFilter
		request.IpFreqLimit = domainResult.IpFreqLimit
		request.StatusCodeCache = domainResult.StatusCodeCache
		request.Compression = domainResult.Compression
		request.BandwidthAlert = domainResult.BandwidthAlert
		request.RangeOriginPull = domainResult.RangeOriginPull
		request.FollowRedirect = domainResult.FollowRedirect
		request.ErrorPage = domainResult.ErrorPage
		request.RequestHeader = domainResult.RequestHeader
		request.ResponseHeader = domainResult.ResponseHeader
		request.DownstreamCapping = domainResult.DownstreamCapping
		request.CacheKey = domainResult.CacheKey
		request.ResponseHeaderCache = domainResult.ResponseHeaderCache
		request.VideoSeek = domainResult.VideoSeek
		request.Cache = domainResult.Cache
		request.OriginPullOptimization = domainResult.OriginPullOptimization
		request.Authentication = domainResult.Authentication
		request.Seo = domainResult.Seo
		request.ForceRedirect = domainResult.ForceRedirect
		request.Referer = domainResult.Referer
		request.MaxAge = domainResult.MaxAge
		request.ServiceType = domainResult.ServiceType
		request.SpecificConfig = domainResult.SpecificConfig
		request.Area = domainResult.Area
		request.OriginPullTimeout = domainResult.OriginPullTimeout
		request.AwsPrivateAccess = domainResult.AwsPrivateAccess
		request.UserAgentFilter = domainResult.UserAgentFilter
		request.AccessControl = domainResult.AccessControl
		request.UrlRedirect = domainResult.UrlRedirect
		request.AccessPort = domainResult.AccessPort
		request.AdvancedAuthentication = domainResult.AdvancedAuthentication
		request.OriginAuthentication = domainResult.OriginAuthentication
		request.Ipv6Access = domainResult.Ipv6Access
		request.OfflineCache = domainResult.OfflineCache
		request.OriginCombine = domainResult.OriginCombine
		request.OssPrivateAccess = domainResult.OssPrivateAccess
		request.WebSocket = domainResult.WebSocket
		request.RemoteAuthentication = domainResult.RemoteAuthentication
		request.Https = domainResult.Https
		if request.Https == nil {
			request.Https = &cdn.Https{
				CertInfo: &cdn.ServerCert{
					CertId: certId,
				},
			}
		}

		_, err := client.UpdateDomainConfig(request)
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			return err
		}
		if err != nil {
			return err
		}

		log.Infof("[%s] Update CDN config for domain %s success", *certId, *domain)

		return nil
	}

	log.Warnf("Not found domain %s", *domain)

	return fmt.Errorf("Not found domain %s", *domain)
}
