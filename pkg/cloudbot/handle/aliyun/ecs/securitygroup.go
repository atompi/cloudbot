package ecs

import (
	"sync"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/atompi/cloudbot/pkg/aliyun/ecs/securitygroup"
	"github.com/atompi/cloudbot/pkg/cloudbot/handle/options"
	"github.com/atompi/cloudbot/pkg/dataio"
	"go.uber.org/zap"
)

func RevokeSecurityGroupHandler(t options.TaskOptions) error {
	res, err := dataio.InputCSV(t.Input)
	if err != nil {
		zap.S().Errorf("input error: %v", err)
		return err
	}

	wg := sync.WaitGroup{}
	ch := make(chan int, t.Threads)

	q := map[string]interface{}{}
	q["RegionId"] = tea.String(t.CloudProvider.RegionId)
	for _, row := range res {
		q["SecurityGroupRuleId.1"] = tea.String(row[0])
		q["SecurityGroupId"] = tea.String(row[1])

		wg.Add(1)
		ch <- 1
		go securitygroup.RevokeSecurityGroup(ch, &wg, t, q)
	}

	wg.Wait()
	return nil
}
