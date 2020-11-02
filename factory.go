package factory

import (
	"fmt"
	"reflect"

	"github.com/vicxu416/gogo-factory/attr"
)

type Template func() interface{}
type Inserter interface {
	Insert(data interface{}) error
}

func New(src Template, attrs ...attr.Attributer) *Factory {
	return &Factory{
		src:          src,
		attrs:        attrs,
		omits:        make(map[string]bool),
		insertQueue:  &ObjectsQueue{q: &Queue{}},
		dependManger: NewDepMan(),
		tempFields:   make([]*tmepField, 0, 1),
		fixFields:    make(map[string]string),
	}
}

type tmepField struct {
	colName string
	field   string
	data    interface{}
}

type Factory struct {
	table         string
	src           Template
	insertFunc    InsertFunc
	attrs         []attr.Attributer
	beforeFactory map[string]*Factory
	afterFactory  map[string]*Factory
	omits         map[string]bool
	tempFields    []*tmepField
	fixFields     map[string]string
	insertQueue   *ObjectsQueue
	dependManger  *DependencyManager
}

func (f *Factory) Table(tableName string) *Factory {
	f.table = tableName
	return f
}

func (f *Factory) Inserter(fn InsertFunc) *Factory {
	f.insertFunc = fn
	return f
}

func (f *Factory) MustBuild() interface{} {
	obj, err := f.build(false)
	f.clear()
	if err != nil {
		panic(err)
	}
	return obj.data
}

func (f *Factory) Build() (interface{}, error) {
	obj, err := f.build(false)
	f.clear()
	if err != nil {
		return nil, err
	}
	return obj.data, nil
}

func (f *Factory) MustInsert() interface{} {
	obj, err := f.build(true)
	if err != nil {
		f.clear()
		panic(err)
	}
	if err := f.insert(); err != nil {
		f.clear()
		panic(err)
	}
	f.clear()
	return obj.data
}

func (f *Factory) Insert() (interface{}, error) {
	obj, err := f.build(true)
	if err != nil {
		f.clear()
		return nil, err
	}
	if err := f.insert(); err != nil {
		f.clear()
		return nil, err
	}
	f.clear()
	return obj.data, nil
}

func (f *Factory) Omit(fields ...string) *Factory {
	for _, field := range fields {
		f.omits[field] = true
	}
	return f
}

func (f *Factory) Fix(fields ...string) *Factory {
	for i := 0; i < len(fields); i += 2 {
		fieldName := fields[i]
		colName := ""
		if i+1 < len(fields) {
			colName = fields[i+1]
		}
		if colName != "" && fieldName != "" {
			f.fixFields[fieldName] = colName
		}
	}
	return f
}

func (f *Factory) FAssociate(name string, other *Factory, n int, before bool, process Processor, options ...string) *Factory {
	depend := &dependency{
		field:   name,
		colName: getColName(options),
		factory: other,
		process: process,
		num:     n,
		fix:     true,
	}
	if before {
		f.dependManger.addBefore(depend)
		return f
	}
	f.dependManger.addAfter(depend)
	return f
}

func (f *Factory) Associate(name string, other *Factory, n int, before bool, process Processor, options ...string) *Factory {
	depend := &dependency{
		field:   name,
		colName: getColName(options),
		factory: other,
		process: process,
		num:     n,
	}
	if before {
		f.dependManger.addBefore(depend)
		return f
	}
	f.dependManger.addAfter(depend)
	return f
}

func (f *Factory) build(insert bool) (*Object, error) {
	data := f.src()
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("build: template data should be a pointer")
	}
	val = val.Elem()

	colVals := make(map[string]interface{})
	for _, attrItem := range f.attrs {
		if f.omits[attrItem.Name()] {
			continue
		}
		field := val.FieldByName(attrItem.Name())
		if !field.IsValid() {
			return nil, fmt.Errorf("build: field(%s) not found", attrItem.Name())
		}
		fieldVal, err := attr.SetField(data, field, attrItem)
		if err != nil {
			return nil, err
		}
		if fieldVal != nil && attrItem.ColName() != "" {
			colVals[attrItem.ColName()] = fieldVal
		}
	}
	f.setFixFields(data, colVals)

	err := f.dependManger.buildBefore(data, f.insertQueue, insert, colVals)
	if err != nil {
		return nil, err
	}

	object := &Object{
		data:       data,
		colVals:    colVals,
		insert:     insert,
		insertFunc: f.insertFunc,
		table:      f.table,
	}
	f.insertQueue.Enqueue(object)
	err = f.dependManger.buildAfter(data, f.insertQueue, insert)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (f *Factory) setFixFields(data interface{}, colVals map[string]interface{}) {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for fieldName, colName := range f.fixFields {
		field := val.FieldByName(fieldName)
		colVals[colName] = field.Interface()
	}
}

func (f *Factory) insert() error {
	object := f.insertQueue.Dequeue()
	var err error
	for object != nil {
		err = object.Insert()
		object = f.insertQueue.Dequeue()
	}
	if err != nil {
		return err
	}
	return nil
}

func (f *Factory) clear() {
	f.insertQueue.clear()
	f.dependManger.clear()
	f.omits = make(map[string]bool)
}
