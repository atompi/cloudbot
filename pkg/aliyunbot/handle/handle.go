package handle

import (
	"github.com/atompi/aliyunbot/pkg/aliyunbot/options"
	"go.uber.org/zap"
)

func Handle(opts options.Options) {
	zap.S().Infof("options: ", opts)
}
