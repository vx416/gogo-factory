package factory

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/vicxu416/seed-factory/attr"
)

type Template func() interface{}

func New(src Template, attrs ...attr.Attributer) *Factory {
	return &Factory{
		src:     src,
		attrs:   attrs,
		colVals: make(map[string]string),
	}
}

type Factory struct {
	table string
	src   Template

	insert  InsertFunc
	attrs   []attr.Attributer
	colVals map[string]string
}

func (f *Factory) TableName(tableName string) *Factory {
	f.table = tableName
	return f
}

func (f *Factory) SetInserter(fn InsertFunc) *Factory {
	f.insert = fn
	return f
}

func (f *Factory) MustBuild() interface{} {
	data, err := f.build(false)
	if err != nil {
		panic(err)
	}
	return data
}

func (f *Factory) Build() (interface{}, error) {
	return f.build(false)
}

func (f *Factory) BuildSeed() (interface{}, error) {
	return f.build(true)
}

func (f *Factory) Insert(data interface{}) error {
	if f.insert != nil {
		return f.insert(options.DB, data)
	}

	if len(f.colVals) == 0 {
		return fmt.Errorf("build and insert: attributes has no column name in %s", f.TableName)
	}

	err := insert(options.DB, f.table, f.colVals)
	if err != nil {
		return err
	}
	return nil
}

func (f *Factory) setAttrValues(fieldVal interface{}, attr attr.Attributer) {
	if attr.ColName() != "" {
		f.colVals[attr.ColName()] = interfaceToStr(fieldVal)
	}
}

func (f *Factory) build(insert bool) (interface{}, error) {
	if insert && options.DB == nil {
		return nil, fmt.Errorf("build with insert: global database instance not exist")
	}

	data := f.src()
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("build: template data should be a pointer")
	}
	val = val.Elem()

	for _, attrItem := range f.attrs {
		field := val.FieldByName(attrItem.Name())
		fieldVal, err := f.setField(data, field, attrItem, insert)
		f.setAttrValues(fieldVal, attrItem)
		if err != nil {
			return nil, err
		}
	}

	if insert {
		err := f.Insert(data)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (f *Factory) setField(data interface{}, field reflect.Value, attrGen attr.Attributer, insert bool) (interface{}, error) {
	val, err := attrGen.Gen(data)
	if err != nil {
		return nil, err
	}

	var fieldRaw interface{}
	if field.CanAddr() {
		fieldRaw = field.Addr().Interface()
	} else {
		fieldRaw = field.Interface()
	}
	if scanner, ok := fieldRaw.(sql.Scanner); ok {
		if err := scanner.Scan(val); err != nil {
			return nil, fmt.Errorf("set scanner field, scan occurs error, %+v", err)
		}
		return val, nil
	}

	if field.Kind() == reflect.Ptr {
		field.Set(reflect.New(field.Type().Elem()))
		field = field.Elem()
	}

	switch attrGen.Kind() {
	case attr.IntAttr:
		realVal := val.(int)
		field.SetInt(int64(realVal))
	case attr.StringAttr:
		realVal := val.(string)
		field.SetString(realVal)
	case attr.FloatAttr:
		realVal := val.(float64)
		field.SetFloat(realVal)
	case attr.BoolAttr:
		realVal := val.(bool)
		field.SetBool(realVal)
	case attr.UnknownAttr:
		field.Set(reflect.ValueOf(val))
	case attr.BytesAttr:
		realVal := val.([]byte)
		field.SetBytes(realVal)
	case attr.UintAttr:
		realVal := val.(uint)
		field.SetUint(uint64(realVal))
	case attr.TimeAttr:
		field.Set(reflect.ValueOf(val))
	case attr.FactoryAttr:
		return f.setFieldWithFactory(data, field, attrGen, insert)
	}

	return val, nil
}

func (f *Factory) setFieldWithFactory(data interface{}, field reflect.Value, factoryAttr attr.Attributer, insert bool) (interface{}, error) {
	factoryData, err := factoryAttr.Gen(data)
	if err != nil {
		return nil, err
	}

	if insert {
		originFactoryAttr := factoryAttr.(attr.Factorier)
		err := originFactoryAttr.Insert(factoryData)
		if err != nil {
			return nil, err
		}
	}
	factoryValue := reflect.ValueOf(factoryData)
	if field.Kind() == reflect.Ptr {
		field.Set(reflect.New(field.Type().Elem()))
		field = field.Elem()
	}

	if factoryValue.Kind() == reflect.Ptr {
		factoryValue = factoryValue.Elem()
	}
	field.Set(factoryValue)
	return factoryData, nil
}
