package handle

import (
	"github.com/alibabacloud-go/tea/tea"
	"github.com/atompi/aliyunbot/pkg/aliyun/slb"
	"github.com/atompi/aliyunbot/pkg/aliyunbot/options"
)

func fetchSLBHandler(t options.TaskOptions) error {
	q := map[string]interface{}{}
	q["RegionId"] = tea.String(t.Aliyun.RegionId)

	err := slb.FetchSLB(t, q)
	return err
}
