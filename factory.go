package factory

import (
	"fmt"
	"reflect"

	"github.com/vicxu416/gogo-factory/attr"
	"github.com/vicxu416/gogo-factory/dbutil"
)

type Template func() reflect.Value

func newTemplate(object interface{}) Template {
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

type ColField func(fieldColumns map[string][]string)

func Col(s ...string) ColField {
	return func(fieldColumns map[string][]string) {
		if len(s) == 2 {
			fieldColumns[s[0]] = []string{s[1]}
		}
		if len(s) == 3 {
			fieldColumns[s[0]] = []string{s[1], s[2]}
		}
	}
}

func New(obj interface{}, attrs ...attr.Attributer) *Factory {
	fieldColumns := make(map[string][]string)

	for _, a := range attrs {
		fieldColumns[a.Name()] = []string{a.ColName()}
	}

	return &Factory{
		src:          newTemplate(obj),
		attrs:        attrs,
		omits:        make(map[string]bool),
		insertQueue:  &ObjectsQueue{q: &Queue{}},
		dependManger: NewDepMan(),
		tempAttr:     make([]attr.Attributer, 0, 1),
		fixFields:    make(map[string]string),
		fieldColumns: fieldColumns,
	}
}

type Factory struct {
	table         string
	src           Template
	insertFunc    dbutil.InsertFunc
	attrs         []attr.Attributer
	beforeFactory map[string]*Factory
	afterFactory  map[string]*Factory
	omits         map[string]bool
	tempAttr      []attr.Attributer
	fixFields     map[string]string
	insertQueue   *ObjectsQueue
	dependManger  *DependencyManager
	fieldColumns  map[string][]string
}

func (f *Factory) Table(tableName string) *Factory {
	f.table = tableName
	return f
}

func (f *Factory) InsertFunc(fn dbutil.InsertFunc) *Factory {
	f.insertFunc = fn
	return f
}

func (f *Factory) MustBuild() interface{} {
	obj, err := f.build(false)
	f.clear()
	if err != nil {
		panic(err)
	}
	return obj.Data
}

func (f *Factory) Build() (interface{}, error) {
	obj, err := f.build(false)
	f.clear()
	if err != nil {
		return nil, err
	}
	return obj.Data, nil
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
	return obj.Data
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
	return obj.Data, nil
}

func (f *Factory) Omit(fields ...string) *Factory {
	for _, field := range fields {
		f.omits[field] = true
	}
	return f
}

func (f *Factory) Columns(fields ...ColField) *Factory {
	for _, field := range fields {
		field(f.fieldColumns)
	}
	return f
}

func (f *Factory) Attrs(attrs ...attr.Attributer) *Factory {
	f.tempAttr = attrs
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

func (f *Factory) build(insert bool) (*dbutil.Object, error) {
	val := f.src()
	data := val.Interface()
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("build: template data should be a pointer")
	}
	val = val.Elem()

	for _, attrItem := range f.attrs {
		if f.omits[attrItem.Name()] {
			continue
		}
		field := val.FieldByName(attrItem.Name())
		if !field.IsValid() {
			return nil, fmt.Errorf("build: field(%s) not found", attrItem.Name())
		}
		fieldType, found := val.Type().FieldByName(attrItem.Name())
		if !found {
			return nil, fmt.Errorf("build: field(%s) not found", attrItem.Name())
		}
		_, err := attr.SetField(val.Interface(), field, fieldType, attrItem)
		if err != nil {
			return nil, err
		}
	}

	for _, attrItem := range f.tempAttr {
		if f.omits[attrItem.Name()] {
			continue
		}
		field := val.FieldByName(attrItem.Name())
		if !field.IsValid() {
			return nil, fmt.Errorf("build: field(%s) not found", attrItem.Name())
		}
		fieldType, found := val.Type().FieldByName(attrItem.Name())
		if !found {
			return nil, fmt.Errorf("build: field(%s) not found", attrItem.Name())
		}
		_, err := attr.SetField(val.Interface(), field, fieldType, attrItem)
		if err != nil {
			return nil, err
		}
	}

	err := f.dependManger.buildBefore(data, f.insertQueue, insert)
	if err != nil {
		return nil, err
	}
	object := &dbutil.Object{
		Data:        data,
		FieldColumn: f.getFieldColumns(),
		NeedInsert:  insert,
		InsertFunc:  f.getInsertFunc(),
		Table:       f.table,
		DB:          options.DB,
		Driver:      options.Driver,
	}
	f.insertQueue.Enqueue(object)
	err = f.dependManger.buildAfter(data, f.insertQueue, insert)

	if err != nil {
		return nil, err
	}
	return object, nil
}

func (f *Factory) getFieldColumns() map[string][]string {
	clonedFieldColumns := make(map[string][]string)
	for k, v := range f.fieldColumns {
		clonedFieldColumns[k] = v
	}
	for _, tempAttr := range f.tempAttr {
		clonedFieldColumns[tempAttr.Name()] = []string{tempAttr.ColName()}
	}
	return clonedFieldColumns
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
	f.tempAttr = make([]attr.Attributer, 0, 1)
	f.omits = make(map[string]bool)
}

func (f *Factory) getInsertFunc() dbutil.InsertFunc {
	if f.insertFunc != nil {
		return f.insertFunc
	}
	return options.InsertFunc
}
