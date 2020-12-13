package gofactory

import (
	"database/sql/driver"
	"reflect"
	"regexp"
	"strings"
)

func getColName(options []string) string {
	if len(options) == 0 {
		return ""
	}
	return options[0]
}

func findField(val reflect.Value, fieldName string) (reflect.Value, reflect.StructField, bool) {
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return reflect.Value{}, reflect.StructField{}, false
	}
	fieldType, found := val.Type().FieldByName(fieldName)
	if !found {
		return reflect.Value{}, reflect.StructField{}, false
	}
	return field, fieldType, true
}

func getElem(data interface{}) reflect.Value {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}

func getFieldElem(val reflect.Value, fieldName string) reflect.Value {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	field := val.FieldByName(fieldName)
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}
	return field
}

func getFieldValue(data interface{}, fieldName string) interface{} {
	val := getElem(data)
	field := val.FieldByName(fieldName)
	if field.IsValid() {
		return field.Interface()
	}
	return nil
}

func makeSlice(data interface{}, cap int) reflect.Value {
	val := reflect.ValueOf(data)
	return reflect.MakeSlice(reflect.SliceOf(val.Type()), 0, cap)
}

func fieldIsSlice(data interface{}, fieldName string) bool {
	field := getFieldElem(getElem(data), fieldName)
	return field.Kind() == reflect.Slice
}

func isPtr(val reflect.Value) bool {
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

// DBTagProcess db tag process
func DBTagProcess(tagGetter TagGetter) string {
	return tagGetter.Get("db")
}

// GormTagProcess gorm tag process
func GormTagProcess(tagGetter TagGetter) string {
	regex := regexp.MustCompile(`column:(.*?( |$))`)
	gormTag := tagGetter.Get("gorm")

	subMatch := regex.FindAllStringSubmatch(gormTag, -1)

	if len(subMatch) == 0 {
		return ""
	}

	firstMatch := subMatch[0]

	if len(firstMatch) < 1 {
		return ""
	}

	return strings.TrimSpace(firstMatch[1])
}
