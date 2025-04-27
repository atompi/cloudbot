package cam

import (
	"github.com/atompi/cloudbot/pkg/cloudbot/handle/options"
	"github.com/atompi/cloudbot/pkg/tencent/cam"
)

func GetCAMUsers(t options.TaskOptions) error {
	// 获取用户列表
	return cam.ListUsers(t)
}
