package dataio

import (
	"fmt"
	"os"

	"github.com/atompi/aliyunbot/pkg/aliyunbot/options"
	"github.com/atompi/aliyunbot/pkg/utils"
	"go.uber.org/zap"
)

func InputCSV(input options.InputOutputOptions) ([][]string, error) {
	if input.Type != "csv" {
		zap.S().Errorf("unknown input type: %v", input.Type)
		return [][]string{}, fmt.Errorf("unknown input type: %v", input.Type)
	}

	f, err := os.Open(input.Path + "/" + input.Target)
	if err != nil {
		zap.S().Errorf("open file error: %v", err)
		return [][]string{}, err
	}
	defer f.Close()

	r := utils.BOMAwareCSVReader(f)
	return r.ReadAll()
}
