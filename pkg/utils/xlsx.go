package utils

import (
	"strings"

	"github.com/xuri/excelize/v2"
)

func XlsxDataToMap(data *[][]string) (records *[]map[string]string, err error) {
	header := []string{}
	records = &[]map[string]string{}
	for i, record := range *data {
		if i == 0 {
			for j := 0; j < len(record); j++ {
				header = append(header, strings.TrimSpace(record[j]))
			}
		} else {
			l := map[string]string{}
			for j := 0; j < len(record); j++ {
				l[header[j]] = strings.TrimSpace(record[j])
			}
			*records = append(*records, l)
		}
	}
	return
}

func GetXlsxRows(f *excelize.File, sheet string) (*[][]string, error) {
	result := [][]string{}

	rows, err := f.Rows(sheet)
	if err != nil {
		return &result, err
	}

	i := 1
	for rows.Next() {
		rowValues := []string{}
		j := 1
		cols, err := f.Cols(sheet)
		if err != nil {
			return &result, err
		}
		for cols.Next() {
			cellName, err := excelize.CoordinatesToCellName(j, i)
			if err != nil {
				return &result, err
			}
			cellValue, err := f.CalcCellValue(sheet, cellName)
			if err != nil {
				return &result, err
			}
			rowValues = append(rowValues, cellValue)
			j += 1
		}
		result = append(result, rowValues)
		i += 1
	}
	return &result, nil
}
