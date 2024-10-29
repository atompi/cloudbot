package rocketmq

import (
	"strconv"
	"sync"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
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

func createApiInfo(action string, pathName string, method string) *openapi.Params {
	return &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String("2022-08-01"),
		Protocol:    tea.String("HTTPS"),
		Method:      tea.String(method),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
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

func CreateTopic(ch chan int, wg *sync.WaitGroup, t options.TaskOptions, instanceId string, topicName string, messageType string, remark string) {
	defer func() { wg.Done(); <-ch }()

	action := "CreateTopic"
	pathName := "/instances/" + instanceId + "/topics/" + topicName
	method := "POST"
	body := map[string]interface{}{
		"messageType": messageType,
		"remark":      remark,
	}

	err := callApi(t, action, pathName, method, body)
	if err != nil {
		zap.S().Errorf("call api failed: %v", err)
		return
	}
}

func CreateConsumerGroup(ch chan int, wg *sync.WaitGroup, t options.TaskOptions, instanceId string, consumerGroupId string, deliveryOrderType string, consumeRetryPolicy string, maxRetryTimes string, deadLetterTargetTopic string, remark string) {
	defer func() { wg.Done(); <-ch }()

	mrt, err := strconv.Atoi(maxRetryTimes)
	if err != nil {
		zap.S().Warnf("convert maxRetryTimes failed: %v, use default: 16", err)
		mrt = 16
	}

	action := "CreateConsumerGroup"
	pathName := "/instances/" + instanceId + "/consumerGroups/" + consumerGroupId
	method := "POST"
	body := map[string]interface{}{
		"deliveryOrderType": deliveryOrderType,
		"consumeRetryPolicy": map[string]interface{}{
			"retryPolicy":           consumeRetryPolicy,
			"maxRetryTimes":         mrt,
			"deadLetterTargetTopic": deadLetterTargetTopic,
		},
		"remark": remark,
	}

	err = callApi(t, action, pathName, method, body)
	if err != nil {
		zap.S().Errorf("call api failed: %v", err)
		return
	}
}

func UpdateConsumerGroup(ch chan int, wg *sync.WaitGroup, t options.TaskOptions, instanceId string, consumerGroupId string, deliveryOrderType string, consumeRetryPolicy string, maxRetryTimes string, deadLetterTargetTopic string, remark string) {
	defer func() { wg.Done(); <-ch }()

	mrt, err := strconv.Atoi(maxRetryTimes)
	if err != nil {
		zap.S().Warnf("convert maxRetryTimes failed: %v, use default: 16", err)
		mrt = 16
	}

	action := "UpdateConsumerGroup"
	pathName := "/instances/" + instanceId + "/consumerGroups/" + consumerGroupId
	method := "PATCH"
	body := map[string]interface{}{
		"deliveryOrderType": deliveryOrderType,
		"consumeRetryPolicy": map[string]interface{}{
			"retryPolicy":           consumeRetryPolicy,
			"maxRetryTimes":         mrt,
			"deadLetterTargetTopic": deadLetterTargetTopic,
		},
		"remark": remark,
	}

	err = callApi(t, action, pathName, method, body)
	if err != nil {
		zap.S().Errorf("call api failed: %v", err)
		return
	}
}
