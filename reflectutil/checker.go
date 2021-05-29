package reflectutil

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
)

func CanSet(val reflect.Value) bool {
	return val.IsValid() && val.CanSet()
}

func IsFieldSlice(data interface{}, fieldName string) bool {
	field := GetFieldElem(GetElem(data), fieldName)
	return field.Kind() == reflect.Slice
}

func IsPtr(val reflect.Value) bool {
	return val.Kind() == reflect.Ptr
}

// IsValuer check the filed implement driver.Valuer or not
func IsValuer(field reflect.Value) (driver.Valuer, bool) {
	var fieldRaw interface{}
	fieldRaw = field.Interface()
	if scanner, ok := fieldRaw.(driver.Valuer); ok {
		return scanner, ok
	}
	if field.CanAddr() {
		fieldRaw = field.Addr().Interface()
	}
	if scanner, ok := fieldRaw.(driver.Valuer); ok {
		return scanner, ok
	}
	return nil, false
}

func TryScan(field reflect.Value, data interface{}) (bool, error) {
	var err error
	scanner, ok := IsScanner(field)
	if ok {
		err = scanner.Scan(data)
	}
	return ok, err
}

func IsScanner(field reflect.Value) (sql.Scanner, bool) {
	var fieldRaw interface{}
	if field.CanAddr() {
		fieldRaw = field.Addr().Interface()
	} else {
		fieldRaw = field.Interface()
	}
	if scanner, ok := fieldRaw.(sql.Scanner); ok {
		return scanner, ok
	}
	return nil, false
}
