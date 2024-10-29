package ons

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

func createApiInfo(action string) *openapi.Params {
	return &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String("2019-02-14"),
		Protocol:    tea.String("HTTPS"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		Pathname:    tea.String("/"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
}

func callApi(t options.TaskOptions, action string, queries map[string]interface{}) error {
	c, err := createApiClient(t.Aliyun)
	if err != nil {
		zap.S().Errorf("create api client failed: %v", err)
		return err
	}

	params := createApiInfo(action)

	runtime := &teautil.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
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
	queries := map[string]interface{}{
		"Topic":       tea.String(topicName),
		"MessageType": tea.String(messageType),
		"Remark":      tea.String(remark),
		"InstanceId":  tea.String(instanceId),
	}

	err := callApi(t, action, queries)
	if err != nil {
		zap.S().Errorf("call api failed: %v", err)
		return
	}
}

func CreateConsumerGroup(ch chan int, wg *sync.WaitGroup, t options.TaskOptions, instanceId string, consumerGroupId string, remark string) {
	defer func() { wg.Done(); <-ch }()

	action := "OnsGroupCreate"
	queries := map[string]interface{}{
		"GroupId":    tea.String(consumerGroupId),
		"Remark":     tea.String(remark),
		"InstanceId": tea.String(instanceId),
	}

	err := callApi(t, action, queries)
	if err != nil {
		zap.S().Errorf("call api failed: %v", err)
		return
	}
}
