package tag

import (
	"fmt"
	"sync"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	teautil "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/atompi/cloudbot/pkg/cloudbot/options"
	"github.com/atompi/cloudbot/pkg/utils"
	"go.uber.org/zap"
)

func createApiClient(opts options.AliyunOptions) (*openapi.Client, error) {
	config := utils.AliyunCreateClientConfig(
		tea.String(opts.AccessKeyId),
		tea.String(opts.AccessKeySecret),
		tea.String(opts.RegionId),
		tea.String(opts.Endpoint),
	)

	return openapi.NewClient(config)
}

func createApiInfo(action string, pathName string, method string) *openapi.Params {
	return &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String("2018-08-28"),
		Protocol:    tea.String("HTTPS"),
		Method:      tea.String(method),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		Pathname:    tea.String(pathName),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
}

func callApi(t options.TaskOptions, action string, pathName string, method string, body map[string]interface{}) error {
	c, err := createApiClient(t.Aliyun)
	if err != nil {
		zap.S().Errorf("create api client failed: %v", err)
		return err
	}

	params := createApiInfo(action, pathName, method)

	runtime := &teautil.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Body: body,
	}

	_, err = c.CallApi(params, request, runtime)
	if err != nil {
		zap.S().Errorf("call api failed: %v", err)
		return err
	}

	return nil
}

func TagResources(ch chan int, wg *sync.WaitGroup, t options.TaskOptions, regionId string, arn string, tags string) {
	defer func() { wg.Done(); <-ch }()

	action := "TagResources"
	pathName := "/"
	method := "POST"
	body := map[string]interface{}{
		"RegionId":      regionId,
		"Tags":          tags,
		"ResourceARN.1": arn,
	}

	err := callApi(t, action, pathName, method, body)
	if err != nil {
		zap.S().Errorf("call api failed: %v", err)
		fmt.Println(arn, tags)
		return
	}
}
