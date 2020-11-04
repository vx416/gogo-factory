package dbutil

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func DefaultInsertFunc(obj *Object) error {
	var (
		colsStr, valuesStr string
		values             = make([]interface{}, 0, 1)
	)
	if obj.DB == nil {
		return fmt.Errorf("insert: global database instance is nil")
	}
	if obj.Table == "" {
		return fmt.Errorf("insert: table name should not be empty")
	}

	for k, v := range obj.ColumnValues() {
		colsStr += k + ", "
		valuesStr += "?, "
		values = append(values, v)
	}
	colsStr = strings.TrimRight(colsStr, ", ")
	valuesStr = strings.TrimRight(valuesStr, ", ")
	insertStmt := "INSERT INTO " + obj.Table + " (" + colsStr + ")" + " VALUES (" + valuesStr + ")"
	insertStmt = rebind(bindType(obj.Driver), insertStmt)
	_, err := obj.DB.Exec(insertStmt, values...)
	if err != nil {
		return fmt.Errorf("sql insert failed, stmt:%s, values:%+v, err:%+v", insertStmt, values, err)
	}
	return nil
}

func GormV2InsertFunc(db *gorm.DB) func(obj *Object) error {
	return func(obj *Object) error {
		return db.Create(obj.Data).Error
	}
}
