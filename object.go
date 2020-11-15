package factory

import (
	"fmt"
	"reflect"

	"github.com/vx416/gogo-factory/attr"
)

type ObjectSetter map[string]attr.Attributer

func (setter ObjectSetter) clone() ObjectSetter {
	newAttr := make(map[string]attr.Attributer)
	for k, v := range setter {
		newAttr[k] = v
	}
	return newAttr
}

func (setter ObjectSetter) SetupObject(val reflect.Value, omits map[string]bool) error {
	data := val.Interface()
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("setup object: object should be a pointer")
	}
	val = val.Elem()

	for _, attrItem := range setter {
		if omits[attrItem.Name()] {
			continue
		}
		field, fieldType, found := findField(val, attrItem.Name())
		if !found {
			return fmt.Errorf("setup object: object field(%s) not found", attrItem.Name())
		}
		_, err := attr.SetField(data, field, fieldType, attrItem)
		if err != nil {
			return err
		}
	}
	return nil
}

type objectConstructor func() reflect.Value

func newConstructor(object interface{}) objectConstructor {
	val := reflect.ValueOf(object)
	objType := val.Type()
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		objType = val.Type()
	}
	return func() reflect.Value {
		obj := reflect.New(objType)
		elem := obj.Elem()
		for i := 0; i < elem.NumField(); i++ {
			elem.Field(i).Set(val.Field(i))
		}
		return obj
	}
}

func getColumnValues(val reflect.Value, fieldColumn map[string]string) map[string]interface{} {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	columnValues := make(map[string]interface{})
	for field, column := range fieldColumn {
		field := val.FieldByName(field)
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		val := field.Interface()
		columnValues[column] = val
	}

	return columnValues
}
