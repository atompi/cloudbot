package handle

import (
	aliyunecs "github.com/atompi/cloudbot/pkg/cloudbot/handle/aliyun/ecs"
	aliyunons "github.com/atompi/cloudbot/pkg/cloudbot/handle/aliyun/ons"
	aliyunrocketmq "github.com/atompi/cloudbot/pkg/cloudbot/handle/aliyun/rocketmq"
	aliyunslb "github.com/atompi/cloudbot/pkg/cloudbot/handle/aliyun/slb"
	aliyuntag "github.com/atompi/cloudbot/pkg/cloudbot/handle/aliyun/tag"
	tencentcam "github.com/atompi/cloudbot/pkg/cloudbot/handle/tencent/cam"
	tencentmonitor "github.com/atompi/cloudbot/pkg/cloudbot/handle/tencent/monitor"
	"github.com/atompi/cloudbot/pkg/cloudbot/options"
	"go.uber.org/zap"
)

func Handle(opts options.Options) {
	for _, t := range opts.Tasks {
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
		case "aliyun_OnsCreateTopic":
			err = aliyunons.CreateTopicHandler(t)
		case "aliyun_OnsCreateConsumerGroup":
			err = aliyunons.CreateConsumerGroupHandler(t)
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
