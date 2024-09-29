package handle

import (
	"sync"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/atompi/aliyunbot/pkg/aliyun/ecs/securitygroup"
	"github.com/atompi/aliyunbot/pkg/aliyunbot/options"
	"github.com/atompi/aliyunbot/pkg/dataio"
	"go.uber.org/zap"
)

func revokeSecurityGroupHandler(t options.TaskOptions) error {
	res, err := dataio.InputCSV(t.Input)
	if err != nil {
		zap.S().Errorf("input error: %v", err)
		return err
	}

	wg := sync.WaitGroup{}
	ch := make(chan int, t.Threads)

	q := map[string]interface{}{}
	q["RegionId"] = tea.String(t.Aliyun.RegionId)
	for _, row := range res {
		q["SecurityGroupRuleId.1"] = tea.String(row[0])
		q["SecurityGroupId"] = tea.String(row[1])

		wg.Add(1)
		ch <- 1
		go securitygroup.RevokeSecurityGroup(ch, &wg, t.Aliyun, q)
	}

	wg.Wait()
	return nil
}
