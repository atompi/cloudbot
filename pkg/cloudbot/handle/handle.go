package handle

import (
	aliyunecs "github.com/atompi/cloudbot/pkg/cloudbot/handle/aliyun/ecs"
	aliyunrocketmq "github.com/atompi/cloudbot/pkg/cloudbot/handle/aliyun/rocketmq"
	aliyunrocketmq4 "github.com/atompi/cloudbot/pkg/cloudbot/handle/aliyun/rocketmq4"
	aliyunslb "github.com/atompi/cloudbot/pkg/cloudbot/handle/aliyun/slb"
	aliyuntag "github.com/atompi/cloudbot/pkg/cloudbot/handle/aliyun/tag"
	"github.com/atompi/cloudbot/pkg/cloudbot/handle/options"
	tencentcam "github.com/atompi/cloudbot/pkg/cloudbot/handle/tencent/cam"
	tencentmonitor "github.com/atompi/cloudbot/pkg/cloudbot/handle/tencent/monitor"
	"github.com/atompi/cloudbot/pkg/utils"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

func LoadTasks(xlsxFile, tasksXlsxSheet string) *[]options.TaskOptions {
	f, err := excelize.OpenFile(xlsxFile)
	if err != nil {
		zap.S().Errorf("open xlsx failed: %v", err)
		return nil
	}
	defer f.Close()

	rows, err := utils.GetXlsxRows(f, tasksXlsxSheet)
	if err != nil {
		zap.S().Errorf("read xlsx failed: %v", err)
		return nil
	}

	data, err := utils.XlsxDataToMap(rows)
	if err != nil {
		zap.S().Errorf("convert xlsx data failed: %v", err)
		return nil
	}

	tasks := []options.TaskOptions{}
	for _, r := range *data {
		t := options.TaskOptions{
			Name:    r["name"],
			Enabled: utils.StringToBool(r["enabled"]),
			Type:    r["type"],
			Threads: utils.StringToInt(r["threads"]),
			CloudProvider: options.CloudProviderOptions{
				AccessKeyId:     r["access_key_id"],
				AccessKeySecret: r["access_key_secret"],
				Endpoint:        r["endpoint"],
				RegionId:        r["region_id"],
			},
			Input: options.InputOutputOptions{
				Path:   r["input.path"],
				Type:   r["input.type"],
				Target: r["input.target"],
			},
			Output: options.InputOutputOptions{
				Path:   r["output.path"],
				Type:   r["output.type"],
				Target: r["output.target"],
			},
		}
		tasks = append(tasks, t)
	}

	return &tasks
}

func Handle(tasks *[]options.TaskOptions) {
	for _, t := range *tasks {
		if !t.Enabled {
			zap.S().Infof("task: %v disabled", t.Name)
			continue
		}

		zap.S().Infof("task: %v started", t.Name)

		var err error

		switch t.Type {
		case "aliyun_RevokeSecurityGroup":
			err = aliyunecs.RevokeSecurityGroupHandler(t)
		case "aliyun_DescribeLoadBalancers":
			err = aliyunslb.FetchSLBHandler(t)
		case "aliyun_RocketMQCreateTopic":
			err = aliyunrocketmq.CreateTopicHandler(t)
		case "aliyun_RocketMQCreateConsumerGroup":
			err = aliyunrocketmq.CreateConsumerGroupHandler(t)
		case "aliyun_RocketMQUpdateConsumerGroup":
			err = aliyunrocketmq.UpdateConsumerGroupHandler(t)
		case "aliyun_RocketMQ4CreateTopic":
			err = aliyunrocketmq4.CreateTopicHandler(t)
		case "aliyun_RocketMQ4CreateConsumerGroup":
			err = aliyunrocketmq4.CreateConsumerGroupHandler(t)
		case "aliyun_TagResources":
			err = aliyuntag.TagResourcesHandler(t)
		case "tencent_GetMonitorData":
			err = tencentmonitor.GetMonitorDataHandler(t)
		case "tencent_GetCAMUsers":
			err = tencentcam.GetCAMUsers(t)
		default:
			zap.S().Warnf("unknown task type: %v", t.Type)
		}

		if err != nil {
			zap.S().Errorf("task: %v failed with error: %v", t.Name, err)
		} else {
			zap.S().Infof("task: %v finished", t.Name)
		}
	}
}
