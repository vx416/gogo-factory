package dbutil

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func DefaultInsertFunc(job *InsertJob) error {
	var (
		colsStr, valuesStr string
		values             = make([]interface{}, 0, 1)
	)
	if job.db == nil {
		return fmt.Errorf("insert: global database instance is nil")
	}
	if job.table == "" {
		return fmt.Errorf("insert: table name should not be empty")
	}

	for k, v := range job.columnValues {
		colsStr += k + ", "
		valuesStr += "?, "
		values = append(values, v)
	}
	colsStr = strings.TrimRight(colsStr, ", ")
	valuesStr = strings.TrimRight(valuesStr, ", ")
	insertStmt := "INSERT INTO " + job.table + " (" + colsStr + ")" + " VALUES (" + valuesStr + ")"
	insertStmt = rebind(bindType(job.driver), insertStmt)
	fmt.Println(insertStmt, values)
	_, err := job.db.Exec(insertStmt, values...)
	if err != nil {
		return fmt.Errorf("sql insert failed, stmt:%s, values:%+v, err:%+v", insertStmt, values, err)
	}
	return nil
}

func GormV2InsertFunc(db *gorm.DB) func(job *InsertJob) error {
	return func(job *InsertJob) error {
		return db.Create(job.GetData()).Error
	}
}
