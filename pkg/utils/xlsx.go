package utils

import "strings"

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
