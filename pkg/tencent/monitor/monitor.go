package monitor

import (
	"encoding/json"

	"github.com/atompi/cloudbot/pkg/cloudbot/options"
	"github.com/atompi/cloudbot/pkg/utils"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

func createApiClient(opts options.TencentOptions, method string, region string) *common.Client {
	credential, cpf := utils.TencentCreateClientConfig(opts.SecretId, opts.SecretKey, opts.Endpoint)
	cpf.HttpProfile.ReqMethod = method
	return common.NewCommonClient(credential, region, cpf)
}

func GetMonitorData(t options.TaskOptions, params string) error {
	client := createApiClient(t.Tencent, "POST", t.Tencent.Region)
	request := tchttp.NewCommonRequest("monitor", "2018-07-24", "GetMonitorData")

	err := request.SetActionParameters(params)
	if err != nil {
		return err
	}

	response := tchttp.NewCommonResponse()

	err = client.Send(request, response)
	if err != nil {
		return err
	}

	data := response.GetBody()

	result := make(map[string]interface{})
	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}

	return nil
}
