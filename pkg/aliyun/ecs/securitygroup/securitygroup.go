package securitygroup

import (
	"sync"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	teautil "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/atompi/cloudbot/pkg/cloudbot/options"
	"github.com/atompi/cloudbot/pkg/utils"
	"go.uber.org/zap"
)

func createApiClient(opts options.AliyunOptions) (*openapi.Client, error) {
	config := utils.CreateClientConfig(
		tea.String(opts.AccessKeyId),
		tea.String(opts.AccessKeySecret),
		tea.String(opts.RegionId),
		tea.String(opts.Endpoint),
	)

	return openapi.NewClient(config)
}

func createApiInfo() *openapi.Params {
	return &openapi.Params{
		Action:      tea.String("RevokeSecurityGroup"),
		Version:     tea.String("2014-05-26"),
		Protocol:    tea.String("HTTPS"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		Pathname:    tea.String("/"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
}

func callApi(a options.AliyunOptions, queries map[string]interface{}) error {
	c, err := createApiClient(a)
	if err != nil {
		return err
	}

	params := createApiInfo()
	runtime := &teautil.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}

	_, err = c.CallApi(params, request, runtime)
	if err != nil {
		return err
	}

	return nil
}

func RevokeSecurityGroup(ch chan int, wg *sync.WaitGroup, t options.TaskOptions, queries map[string]interface{}) {
	defer func() { wg.Done(); <-ch }()

	err := callApi(t.Aliyun, queries)
	if err != nil {
		zap.S().Errorf("call api failed: %v", err)
		return
	}
}
