package gofactory

import (
	"fmt"
	"reflect"

	"github.com/vx416/gogo-factory/attr"
	"github.com/vx416/gogo-factory/reflectutil"
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

func (setter ObjectSetter) buildFieldColumns(obj interface{}) map[string]string {
	fieldColumns := make(map[string]string)

	for _, a := range setter {
		fieldColumns[a.Name()] = a.ColName()
	}

	return fieldColumns
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
	field, fieldType, found := reflectutil.FindField(val, attrItem.Name())
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
			if !val.Field(i).IsZero() {
				elem.Field(i).Set(val.Field(i))
			}
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
		if !field.IsZero() {
			val := field.Interface()
			columnValues[column] = val
		}
	}

	return columnValues
}

// TagGetter tag getter
type TagGetter interface {
	Get(tagName string) (tagString string)
}

type TagProcess func(tagGetter TagGetter) string

func getObjectColumnNames(val reflect.Value, tagProcess TagProcess) map[string]string {
	objType := val.Type()
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		objType = val.Type()
	}

	filedColumns := make(map[string]string)

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldName := field.Name
		colName := tagProcess(field.Tag)
		if colName != "" {
			filedColumns[fieldName] = colName
		}
	}

	return filedColumns
}
