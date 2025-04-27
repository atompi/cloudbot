package cam

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/atompi/cloudbot/pkg/cloudbot/handle/options"
	"github.com/atompi/cloudbot/pkg/utils"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	"go.uber.org/zap"
)

func createApiClient(opts options.CloudProviderOptions, method string, region string) *common.Client {
	credential, cpf := utils.TencentCreateClientConfig(opts.AccessKeyId, opts.AccessKeySecret, opts.Endpoint)
	cpf.HttpProfile.ReqMethod = method
	return common.NewCommonClient(credential, region, cpf)
}

func ListUsers(t options.TaskOptions) (err error) {
	client := createApiClient(t.CloudProvider, "POST", t.CloudProvider.RegionId)
	request := tchttp.NewCommonRequest("cam", "2019-01-16", "ListUsers")
	params := "{}"

	err = request.SetActionParameters(params)
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
	users, _ := resp["Data"].([]interface{})
	for _, user := range users {
		time.Sleep(500 * time.Millisecond)
		user, _ := user.(map[string]interface{})
		nickName, _ := user["NickName"].(string)
		name, _ := user["Name"].(string)
		remark, _ := user["Remark"].(string)
		uin, _ := user["Uin"].(float64)
		consoleLogin := user["ConsoleLogin"].(float64)

		request = tchttp.NewCommonRequest("cam", "2019-01-16", "ListAttachedUserAllPolicies")
		params = `{
    "TargetUin": %v,
    "Rp": 200,
    "Page": 1,
    "AttachType": 0
}`
		params = fmt.Sprintf(params, uin)

		err = request.SetActionParameters(params)
		if err != nil {
			zap.S().Errorf("failed to set action parameters: %v", err)
			return
		}
		response = tchttp.NewCommonResponse()
		err = client.Send(request, response)
		if err != nil {
			zap.S().Errorf("failed to send request: %v", err)
			return
		}
		data = response.GetBody()
		result = make(map[string]interface{})
		err = json.Unmarshal(data, &result)
		if err != nil {
			zap.S().Errorf("failed to unmarshal json: %v", err)
			return
		}
		resp, _ = result["Response"].(map[string]interface{})
		policies, _ := resp["PolicyList"].([]interface{})
		for _, policy := range policies {
			policy, _ := policy.(map[string]interface{})
			policyName, _ := policy["PolicyName"].(string)
			description, _ := policy["Description"].(string)
			deactived, _ := policy["Deactived"].(bool)
			fmt.Printf("%v,%v,%v,%v,%v,%v,%v\n", name, nickName, remark, consoleLogin, policyName, description, deactived)
		}
	}
	return
}
