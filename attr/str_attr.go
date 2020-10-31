package attr

import (
	"fmt"
)

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

func (strAttr) Kind() AttrType {
	return StringAttr
}

func (attr *strAttr) Gen(data interface{}) (interface{}, error) {
	attr.val = attr.genFunc()
	if attr.process != nil {
		if err := attr.process(attr, data); err != nil {
			return nil, err
		}
	}
	return attr.val, nil
}

func (attr strAttr) Name() string {
	return attr.name
}

func StrSeq(name string, strSeq []string, options ...string) Attributer {
	return &strSeqAttr{
		name:    name,
		colName: getColName(options),
		strSeq:  strSeq,
		index:   0,
	}
}

type strSeqAttr struct {
	name    string
	colName string
	strSeq  []string
	index   int
	val     string
	process Processor
}

func (attr *strSeqAttr) Process(procFunc Processor) Attributer {
	attr.process = procFunc
	return attr
}

func (attr strSeqAttr) GetVal() interface{} {
	return attr.val
}

func (attr *strSeqAttr) SetVal(val interface{}) error {
	realVal, ok := val.(string)
	if !ok {
		return fmt.Errorf("set attribute val: val %+v is not string", val)
	}

	attr.val = realVal
	return nil
}

func (attr strSeqAttr) ColName() string {
	return attr.colName
}

func (strSeqAttr) Kind() AttrType {
	return StringAttr
}

func (attr *strSeqAttr) Gen(data interface{}) (interface{}, error) {
	attr.val = attr.strSeq[attr.index]
	if attr.process != nil {
		if err := attr.process(attr, data); err != nil {
			return nil, err
		}
	}

	attr.index++
	if attr.index >= len(attr.strSeq) {
		attr.index = 0
	}
	return attr.val, nil
}

func (attr strSeqAttr) Name() string {
	return attr.name
}
