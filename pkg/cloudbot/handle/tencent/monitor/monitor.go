package monitor

import (
	"encoding/json"

	"github.com/atompi/cloudbot/pkg/cloudbot/options"
	"github.com/atompi/cloudbot/pkg/dataio"
	"github.com/atompi/cloudbot/pkg/tencent/monitor"
	"github.com/atompi/cloudbot/pkg/utils"
	"go.uber.org/zap"
)

func GetMonitorDataHandler(t options.TaskOptions) error {
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

	for _, row := range *data {
		paramsMap := map[string]interface{}{
			"Namespace":  row["namespace"],
			"MetricName": row["mem"],
			"Period":     86400,
			"StartTime":  "2024-09-29T00:00:00+08:00",
			"EndTime":    "2024-10-29T00:00:00+08:00",
			"Instances": []map[string][]map[string]string{
				{
					"Dimensions": []map[string]string{
						{
							"Name":  row["name"],
							"Value": row["value"],
						},
					},
				},
			},
			"SpecifyStatistics": 2,
		}
		params, _ := json.Marshal(paramsMap)

		monitor.GetMonitorData(t, string(params), row["value"])
	}

	return nil
}
