package monitor

import (
	"encoding/json"
	"fmt"

	"github.com/atompi/cloudbot/pkg/cloudbot/options"
	"github.com/atompi/cloudbot/pkg/utils"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	"go.uber.org/zap"
)

func createApiClient(opts options.TencentOptions, method string, region string) *common.Client {
	credential, cpf := utils.TencentCreateClientConfig(opts.SecretId, opts.SecretKey, opts.Endpoint)
	cpf.HttpProfile.ReqMethod = method
	return common.NewCommonClient(credential, region, cpf)
}

func GetMonitorData(t options.TaskOptions, params string, instanceId string) {
	client := createApiClient(t.Tencent, "POST", t.Tencent.Region)
	request := tchttp.NewCommonRequest("monitor", "2018-07-24", "GetMonitorData")

	err := request.SetActionParameters(params)
	if err != nil {
		zap.S().Errorf("failed to set action parameters: %v", err)
		return
	}

	response := tchttp.NewCommonResponse()

	err = client.Send(request, response)
	if err != nil {
		zap.S().Errorf("failed to send request: %v", err)
		return
	}

	data := response.GetBody()

	result := make(map[string]interface{})
	err = json.Unmarshal(data, &result)
	if err != nil {
		zap.S().Errorf("failed to unmarshal json: %v", err)
		return
	}

	resp, _ := result["Response"].(map[string]interface{})
	dataPoints, _ := resp["DataPoints"].([]interface{})
	dataPoint, _ := dataPoints[0].(map[string]interface{})
	maxValuesInterface, _ := dataPoint["MaxValues"].([]interface{})
	maxValues := []float64{}
	for _, v := range maxValuesInterface {
		maxValues = append(maxValues, v.(float64))
	}
	if len(maxValues) == 0 {
		return
	}
	maxValue := utils.MaxFloat64(maxValues)

	fmt.Println(instanceId, maxValue)
}
