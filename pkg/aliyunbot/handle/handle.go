package handle

import (
	"github.com/atompi/aliyunbot/pkg/aliyunbot/handle/ecs"
	"github.com/atompi/aliyunbot/pkg/aliyunbot/handle/rocketmq"
	"github.com/atompi/aliyunbot/pkg/aliyunbot/handle/slb"
	"github.com/atompi/aliyunbot/pkg/aliyunbot/options"
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
		case "RevokeSecurityGroup":
			err = ecs.RevokeSecurityGroupHandler(t)
		case "DescribeLoadBalancers":
			err = slb.FetchSLBHandler(t)
		case "RocketMQCreateTopic":
			err = rocketmq.CreateTopicHandler(t)
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
