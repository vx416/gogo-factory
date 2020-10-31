package attr

import (
	"reflect"
)

func Int(name string, genFunc func() int, options ...string) Attributer {
	return &intAttribute{
		name:    name,
		colName: options[0],
		genFunc: func(data interface{}) (int, error) {
			return genFunc(), nil
		},
	}
}

func IntWith(name string, genFunc func(data interface{}) (int, error), options ...string) Attributer {
	return &intAttribute{
		name:    name,
		colName: options[0],
		genFunc: genFunc,
	}
}

type intAttribute struct {
	name    string
	colName string
	genFunc func(data interface{}) (int, error)
}

func (attr intAttribute) ColName() string {
	return attr.colName
}

func (attr intAttribute) Name() string {
	return attr.name
}

func (intAttribute) Kind() reflect.Kind {
	return reflect.Int64
}

func (attr intAttribute) Gen(data interface{}) (interface{}, error) {
	return attr.genFunc(data)
}

func Seq(name string, start int, options ...string) Attributer {
	return &seqAttr{
		name:    name,
		colName: options[0],
		seq:     start,
		genFunc: func(data interface{}, seq int) (int, error) {
			return seq, nil
		},
	}
}

func SeqWith(name string, genFunc func(data interface{}, seq int) (int, error), options ...string) Attributer {
	return &seqAttr{
		name:    name,
		colName: options[0],
		genFunc: genFunc,
		seq:     1}
}

type seqAttr struct {
	seq     int
	name    string
	colName string
	genFunc func(data interface{}, seq int) (int, error)
}

func (attr seqAttr) ColName() string {
	return attr.colName
}

func (attr *seqAttr) Gen(data interface{}) (interface{}, error) {
	val, err := attr.genFunc(data, attr.seq)
	if err != nil {
		return nil, err
	}
	attr.seq++
	return val, nil
}

func (seqAttr) Kind() reflect.Kind {
	return reflect.Int64
}

func (attr seqAttr) Name() string {
	return attr.name
}

func Float(name string, genFunc func() float64, options ...string) Attributer {
	return &floatAttr{
		name:    name,
		colName: options[0],
		genFunc: func(data interface{}) (float64, error) {
			return genFunc(), nil
		},
	}
}

func FloatWith(name string, genFunc func(data interface{}) (float64, error), options ...string) Attributer {
	return &floatAttr{
		name:    name,
		colName: options[0],
		genFunc: genFunc,
	}
}

type floatAttr struct {
	name    string
	colName string
	genFunc func(data interface{}) (float64, error)
}

func (attr floatAttr) ColName() string {
	return attr.colName
}

func (attr floatAttr) Gen(data interface{}) (interface{}, error) {
	return attr.genFunc(data)
}

func (floatAttr) Kind() reflect.Kind {
	return reflect.Float64
}

func (attr floatAttr) Name() string {
	return attr.name
}

func Uint(name string, genFunc func() uint, options ...string) Attributer {
	return &uintAttr{
		name:    name,
		colName: options[0],
		genFunc: func(data interface{}) (uint, error) {
			return genFunc(), nil
		},
	}
}

func UintWith(name string, genFunc func(data interface{}) (uint, error), options ...string) Attributer {
	return &uintAttr{
		name:    name,
		colName: options[0],
		genFunc: genFunc,
	}
}

type uintAttr struct {
	name    string
	colName string
	genFunc func(data interface{}) (uint, error)
}

func (attr uintAttr) ColName() string {
	return attr.colName
}

func (attr uintAttr) Gen(data interface{}) (interface{}, error) {
	return attr.genFunc(data)
}

func (uintAttr) Kind() reflect.Kind {
	return reflect.Uint
}

func (attr uintAttr) Name() string {
	return attr.name
}
