package attr

import (
	"fmt"
)

// Str create string attributer with generated function
func Str(name string, genFunc func() string, options ...string) Attributer {
	return &strAttr{
		name:    name,
		colName: getColName(options),
		genFunc: genFunc,
	}
}

type strAttr struct {
	val     string
	name    string
	colName string
	genFunc func() string
	process Processor
	object  interface{}
}

func (attr *strAttr) GetObject() interface{} {
	return attr.object
}

func (attr *strAttr) Process(procFunc Processor) Attributer {
	attr.process = procFunc
	return attr
}

func (attr strAttr) GetVal() interface{} {
	return attr.val
}

func (attr *strAttr) SetVal(val interface{}) error {
	realVal, ok := val.(string)
	if !ok {
		return fmt.Errorf("set attribute val: val %+v is not string", val)
	}

	attr.val = realVal
	return nil
}

func (attr strAttr) ColName() string {
	return attr.colName
}

func (strAttr) Kind() Type {
	return StringAttr
}

func (attr *strAttr) Gen(data interface{}) (interface{}, error) {
	attr.val = attr.genFunc()
	attr.object = data
	if attr.process != nil {
		if err := attr.process(attr); err != nil {
			return nil, err
		}
	}
	return attr.val, nil
}

func (attr strAttr) Name() string {
	return attr.name
}
