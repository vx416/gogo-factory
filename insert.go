package factory

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	UNKNOWN = iota
	QUESTION
	DOLLAR
	NAMED
	AT
)

type InsertFunc func(db *sql.DB, data interface{}) error

func bindType(driverName string) int {
	switch driverName {
	case "postgres", "pgx", "pq-timeouts", "cloudsqlpostgres", "ql", "pg":
		return DOLLAR
	case "mysql":
		return QUESTION
	case "sqlite3", "sqlite":
		return QUESTION
	case "oci8", "ora", "goracle", "godror":
		return NAMED
	case "sqlserver":
		return AT
	}
	return UNKNOWN
}

func rebind(bindType int, query string) string {
	switch bindType {
	case QUESTION, UNKNOWN:
		return query
	}

	// Add space enough for 10 params before we have to allocate
	rqb := make([]byte, 0, len(query)+10)

	var i, j int

	for i = strings.Index(query, "?"); i != -1; i = strings.Index(query, "?") {
		rqb = append(rqb, query[:i]...)

		switch bindType {
		case DOLLAR:
			rqb = append(rqb, '$')
		case NAMED:
			rqb = append(rqb, ':', 'a', 'r', 'g')
		case AT:
			rqb = append(rqb, '@', 'p')
		}

		j++
		rqb = strconv.AppendInt(rqb, int64(j), 10)

		query = query[i+1:]
	}

	return string(append(rqb, query...))
}

func insert(db *sql.DB, tableName string, colVal map[string]interface{}) error {
	var (
		colsStr, valuesStr string
		values             = make([]interface{}, 0, len(colVal))
	)
	if db == nil {
		return fmt.Errorf("insert: global database instance is nil")
	}
	if tableName == "" {
		return fmt.Errorf("insert: table name should not be empty")
	}

	for k, v := range colVal {
		colsStr += k + ", "
		valuesStr += "?, "
		values = append(values, v)
	}
	colsStr = strings.TrimRight(colsStr, ", ")
	valuesStr = strings.TrimRight(valuesStr, ", ")
	insertStmt := "INSERT INTO " + tableName + " (" + colsStr + ")" + " VALUES (" + valuesStr + ")"
	insertStmt = rebind(bindType(options.driver), insertStmt)
	_, err := db.Exec(insertStmt, values...)
	if err != nil {
		return fmt.Errorf("sql insert failed, stmt:%s, err:%+v", err, insertStmt)
	}

	return nil
}

func getID(in interface{}) interface{} {
	val := reflect.ValueOf(in)

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	id := val.FieldByName("ID")
	if !id.IsValid() {
		id = val.FieldByName("Id")
	}
	if !id.IsValid() {
		return nil
	}
	return id.Interface()
}
