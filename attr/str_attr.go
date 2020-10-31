package attr

import (
	"reflect"
)

func Str(name string, genFunc func() string, options ...string) Attributer {
	return &strAttr{
		name:    name,
		colName: options[0],
		genFunc: func(data interface{}) (string, error) {
			return genFunc(), nil
		},
	}
}

func StrWith(name string, genFunc func(data interface{}) (string, error), options ...string) Attributer {
	return &strAttr{
		name:    name,
		colName: options[0],
		genFunc: genFunc,
	}
}

type strAttr struct {
	name    string
	colName string
	genFunc func(data interface{}) (string, error)
}

func (attr strAttr) ColName() string {
	return attr.colName
}

func (strAttr) Kind() reflect.Kind {
	return reflect.String
}

func (attr strAttr) Gen(data interface{}) (interface{}, error) {
	return attr.genFunc(data)
}

func (attr strAttr) Name() string {
	return attr.name
}

func StrSeq(name string, strSeq []string, options ...string) Attributer {
	return &strSeqAttr{
		name:    name,
		colName: options[0],
		strSeq:  strSeq,
		index:   0,
		genFunc: func(data interface{}, str string) (string, error) {
			return str, nil
		},
	}
}

func StrSeqWith(name string, strSeq []string, genFunc func(data interface{}, str string) (string, error), options ...string) Attributer {
	return &strSeqAttr{
		name:    name,
		colName: options[0],
		strSeq:  strSeq,
		index:   0,
		genFunc: genFunc,
	}
}

type strSeqAttr struct {
	name    string
	colName string
	strSeq  []string
	index   int
	genFunc func(data interface{}, str string) (string, error)
}

func (attr strSeqAttr) ColName() string {
	return attr.colName
}

func (strSeqAttr) Kind() reflect.Kind {
	return reflect.String
}

func (attr *strSeqAttr) Gen(data interface{}) (interface{}, error) {
	str, err := attr.genFunc(data, attr.strSeq[attr.index])
	if err != nil {
		return nil, err
	}

	attr.index++
	if attr.index >= len(attr.strSeq) {
		attr.index = 0
	}
	return str, nil
}

func (attr strSeqAttr) Name() string {
	return attr.name
}
