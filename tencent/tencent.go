package tencent

import (
	"github.com/neobaran/csac/config"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

type TencentCloudHelper struct {
	credential *common.Credential
}

func NewTencentCloudHelp(config config.TencentCredential) *TencentCloudHelper {
	return &TencentCloudHelper{
		credential: common.NewCredential(
			config.SecretId,
			config.SecretKey,
		),
	}
}
