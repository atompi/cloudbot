package handle

import (
	"os"
	"sync"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/atompi/aliyunbot/pkg/aliyun/ecs/securitygroup"
	"github.com/atompi/aliyunbot/pkg/aliyunbot/options"
	"github.com/atompi/aliyunbot/pkg/utils"
	"go.uber.org/zap"
)

func Handle(opts options.Options) {
	for _, t := range opts.Tasks {
		zap.S().Infof("start handle task: %v", t.Name)

		f, err := os.Open(t.InputFile)
		if err != nil {
			zap.S().Errorf("open file error: %v", err)
			return
		}
		defer f.Close()

		r := utils.BOMAwareCSVReader(f)
		res, err := r.ReadAll()
		if err != nil {
			zap.S().Errorf("read csv error: %v", err)
			return
		}

		wg := sync.WaitGroup{}
		ch := make(chan int, opts.Core.Threads)

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
	}
}
