package attr

import (
	"reflect"
)

type Attributer interface {
	Name() string
	ColName() string
	Kind() reflect.Kind
	Gen(data interface{}) (interface{}, error)
}

type Factorier interface {
	Build() (interface{}, error)
	Insert(data interface{}) error
}

func Attr(name string, genFunc func() interface{}, options ...string) Attributer {
	return &attr{
		name:    name,
		colName: options[0],
		genFunc: func(data interface{}) (interface{}, error) {
			return genFunc(), nil
		},
	}
}

func AttrWith(name string, genFunc func(data interface{}) (interface{}, error), options ...string) Attributer {
	return &attr{
		name:    name,
		colName: options[0],
		genFunc: genFunc,
	}
}

type attr struct {
	name    string
	colName string
	genFunc func(data interface{}) (interface{}, error)
}

func (attr attr) ColName() string {
	return attr.colName
}

func (attr attr) Name() string {
	return attr.name
}

func (attr) Kind() reflect.Kind {
	return reflect.Interface
}

func (attr attr) Gen(data interface{}) (interface{}, error) {
	return attr.genFunc(data)
}
