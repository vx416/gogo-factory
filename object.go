package factory

import (
	"fmt"
	"reflect"

	"github.com/vx416/gogo-factory/attr"
)

// ObjectSetter object setter
type ObjectSetter []attr.Attributer

func (setter ObjectSetter) clone() ObjectSetter {
	newAttr := make([]attr.Attributer, len(setter))
	for i, v := range setter {
		newAttr[i] = v
	}
	return newAttr
}

// SetupObject setup object with Attributers
func (setter ObjectSetter) SetupObject(val reflect.Value, omits map[string]bool, only map[string]bool) error {
	data := val.Interface()
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("setup object: object should be a pointer")
	}
	val = val.Elem()

	if len(only) > 0 {
		for _, attrItem := range setter {
			if !only[attrItem.Name()] {
				continue
			}
			if omits[attrItem.Name()] {
				continue
			}
			err := setter.setField(data, val, attrItem)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, attrItem := range setter {
		if omits[attrItem.Name()] {
			continue
		}
		err := setter.setField(data, val, attrItem)
		if err != nil {
			return err
		}
	}
	return nil
}

func (setter ObjectSetter) setField(data interface{}, val reflect.Value, attrItem attr.Attributer) error {
	field, fieldType, found := findField(val, attrItem.Name())
	if !found {
		return fmt.Errorf("setup object: object field(%s) not found", attrItem.Name())
	}
	_, err := attr.SetField(data, field, fieldType, attrItem)
	if err != nil {
		return err
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
