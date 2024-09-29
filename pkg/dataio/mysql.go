package dataio

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/atompi/aliyunbot/pkg/aliyunbot/options"
	"go.uber.org/zap"
)

func OutputMySQL(fields []string, values []interface{}, output options.InputOutputOptions) error {
	if output.Type != "mysql" {
		zap.S().Errorf("unknown input type: %v", output.Type)
		return fmt.Errorf("unknown input type: %v", output.Type)
	}

	db, err := sql.Open("mysql", output.Path)
	if err != nil {
		return fmt.Errorf("error opening database connection: %v", err)
	}
	defer db.Close()

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		output.Target,
		strings.Join(fields, ", "),
		strings.Repeat("?, ", len(fields)-1)+"?",
	)

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare statement failed: %v", err)
	}
	defer stmt.Close()

	if len(values) != len(fields) {
		return fmt.Errorf("number of values and fields do not match")
	}
	_, err = stmt.Exec(values...)
	if err != nil {
		return fmt.Errorf("executing statement failed: %w", err)
	}

	return nil
}
