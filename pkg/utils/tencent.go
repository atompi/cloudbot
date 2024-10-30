package utils

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

func TencentCreateClientConfig(secretId, secretKey, endpoint string) (*common.Credential, *profile.ClientProfile) {
	credential := common.NewCredential(secretId, secretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = endpoint
	return credential, cpf
}
