package slb

import (
	"github.com/alibabacloud-go/tea/tea"
	"github.com/atompi/cloudbot/pkg/aliyun/slb"
	"github.com/atompi/cloudbot/pkg/cloudbot/handle/options"
)

func FetchSLBHandler(t options.TaskOptions) error {
	q := map[string]interface{}{}
	q["RegionId"] = tea.String(t.CloudProvider.RegionId)

	err := slb.FetchSLB(t, q)
	return err
}
