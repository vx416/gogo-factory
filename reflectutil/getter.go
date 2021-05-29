package reflectutil

import "reflect"

func GetFieldElem(val reflect.Value, fieldName string) reflect.Value {
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

func GetFieldValue(data interface{}, fieldName string) interface{} {
	val := GetElem(data)
	field := val.FieldByName(fieldName)
	if field.IsValid() {
		return field.Interface()
	}
	return nil
}

func GetElem(data interface{}) reflect.Value {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}

func FindField(val reflect.Value, fieldName string) (reflect.Value, reflect.StructField, bool) {
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

func MakeSlice(data interface{}, cap int) reflect.Value {
	val := reflect.ValueOf(data)
	return reflect.MakeSlice(reflect.SliceOf(val.Type()), 0, cap)
}
