package utils

import (
	"encoding/csv"
	"io"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

const UTF8BOM string = "\xEF\xBB\xBF"

func BOMAwareCSVReader(reader io.Reader) *csv.Reader {
	transformer := unicode.BOMOverride(encoding.Nop.NewDecoder())
	return csv.NewReader(transform.NewReader(reader, transformer))
}

func DataToMap(data *[][]string) (records *[]map[string]string, err error) {
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
