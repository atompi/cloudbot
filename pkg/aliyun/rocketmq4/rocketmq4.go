package rocketmq4

import (
	"sync"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	teautil "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/atompi/cloudbot/pkg/cloudbot/handle/options"
	"github.com/atompi/cloudbot/pkg/utils"
	"go.uber.org/zap"
)

func createApiClient(opts options.CloudProviderOptions) (*openapi.Client, error) {
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
		Version:     tea.String("2019-02-14"),
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
	c, err := createApiClient(t.CloudProvider)
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

func CreateTopic(ch chan int, wg *sync.WaitGroup, t options.TaskOptions, instanceId string, topicName string, messageType string, remark string) {
	defer func() { wg.Done(); <-ch }()

	action := "OnsTopicCreate"
	pathName := "/"
	method := "POST"
	body := map[string]interface{}{
		"Topic":       tea.String(topicName),
		"MessageType": tea.String(messageType),
		"Remark":      tea.String(remark),
		"InstanceId":  tea.String(instanceId),
	}

	err := callApi(t, action, pathName, method, body)
	if err != nil {
		zap.S().Errorf("call api failed: %v", err)
		return
	}
}

func CreateConsumerGroup(ch chan int, wg *sync.WaitGroup, t options.TaskOptions, instanceId string, groupId string, groupType string, remark string) {
	defer func() { wg.Done(); <-ch }()

	action := "OnsGroupCreate"
	pathName := "/"
	method := "POST"
	body := map[string]interface{}{
		"GroupId":    groupId,
		"Remark":     remark,
		"InstanceId": instanceId,
		"GroupType":  groupType,
	}

	err := callApi(t, action, pathName, method, body)
	if err != nil {
		zap.S().Errorf("call api failed: %v", err)
		return
	}
}
