package factory

import (
	"reflect"
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
