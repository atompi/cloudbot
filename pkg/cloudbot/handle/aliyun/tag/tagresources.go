package tag

import (
	"fmt"
	"strings"
	"sync"

	"github.com/atompi/cloudbot/pkg/aliyun/tag"
	"github.com/atompi/cloudbot/pkg/cloudbot/options"
	"github.com/atompi/cloudbot/pkg/dataio"
	"github.com/atompi/cloudbot/pkg/utils"
	"go.uber.org/zap"
)

func TagResourcesHandler(t options.TaskOptions) error {
	res, err := dataio.InputCSV(t.Input)
	if err != nil {
		zap.S().Errorf("input error: %v", err)
		return err
	}

	data, err := utils.DataToMap(&res)
	if err != nil {
		zap.S().Errorf("data convert error: %v", err)
		return err
	}

	wg := sync.WaitGroup{}
	ch := make(chan int, t.Threads)

	for _, row := range *data {
		wg.Add(1)
		ch <- 1

		regionId := row["regionId"]
		if regionId == "global" {
			regionId = "cn-hangzhou"
		}
		tags := fmt.Sprintf("{\"%s\":\"%s\"}", row["tagKey"], row["tagValue"])
		a := strings.Split(row["resourceTypeCode"], "::")
		arnPrefix := strings.ToLower("arn:" + a[0] + ":" + a[1])
		arnResourceType := strings.ToLower(a[2])
		arnAccount := row["accountId"]
		arnInstanceId := row["resourceId"]
		arn := fmt.Sprintf("%s:%s:%s:%s/%s", arnPrefix, regionId, arnAccount, arnResourceType, arnInstanceId)

		// fmt.Println(arn, tags)
		go tag.TagResources(ch, &wg, t, regionId, arn, tags)
	}

	wg.Wait()
	return nil
}
