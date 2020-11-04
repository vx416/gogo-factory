package dbutil

import (
	"database/sql"
	"reflect"
)

type InsertFunc func(obj *Object) error

type Object struct {
	Data         interface{}
	FieldColumn  map[string][]string
	DB           *sql.DB
	Driver       string
	Table        string
	NeedInsert   bool
	InsertFunc   InsertFunc
	columnValues map[string]interface{}
	val          reflect.Value
	tag          string
}

func (obj *Object) Insert() error {
	if !obj.NeedInsert {
		return nil
	}
	if obj.InsertFunc != nil {
		return obj.InsertFunc(obj)
	}
	return DefaultInsertFunc(obj)
}

func (obj *Object) ColumnValues() map[string]interface{} {
	if len(obj.columnValues) == 0 {
		obj.columnValues = make(map[string]interface{})
		val := obj.objVal()
		for k, v := range obj.FieldColumn {
			field := val.FieldByName(k)
			if field.Kind() == reflect.Ptr {
				field = field.Elem()
			}
			val := field.Interface()
			column := v[0]
			if len(v) == 2 {
				structField := v[0]
				column = v[1]
				structFieldVal := field.FieldByName(structField)
				if structFieldVal.IsValid() && !structFieldVal.IsZero() {
					val = structFieldVal.Interface()
				}
			}
			obj.columnValues[column] = val
		}
	}
	return obj.columnValues
}

func (obj *Object) objVal() reflect.Value {
	val := reflect.ValueOf(obj.Data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	obj.val = val

	return obj.val
}
