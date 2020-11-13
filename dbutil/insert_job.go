package dbutil

import (
	"database/sql"
	"reflect"
)

type InsertFunc func(job *InsertJob) error

func NewJob(val reflect.Value, columnValues map[string]interface{}) *InsertJob {
	return &InsertJob{
		val:          val,
		columnValues: columnValues,
	}
}

type InsertJob struct {
	db           *sql.DB
	driver       string
	table        string
	insertFunc   InsertFunc
	columnValues map[string]interface{}
	val          reflect.Value
	tag          string
}

func (job *InsertJob) SetDB(db *sql.DB, driver, table, tag string) {
	job.db = db
	job.driver = driver
	job.table = table
	job.tag = tag
}

func (job *InsertJob) SetInsertFunc(fn InsertFunc) {
	job.insertFunc = fn
}

func (job *InsertJob) Insert() error {
	if job.insertFunc != nil {
		return job.insertFunc(job)
	}
	return DefaultInsertFunc(job)
}

func (job *InsertJob) jobVal() reflect.Value {
	val := job.val
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val
}

func (job *InsertJob) GetData() interface{} {
	return job.val.Interface()
}
