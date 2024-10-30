package utils

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
)

func AliyunCreateClientConfig(accessKeyId, accessKeySecret, regionId, endpoint *string) *openapi.Config {
	return &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		RegionId:        regionId,
		Endpoint:        endpoint,
	}
}
