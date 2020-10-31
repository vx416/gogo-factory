package attr

import (
	"fmt"
	"time"
)

func Bool(name string, genFunc func() bool, options ...string) Attributer {
	return &boolAttr{
		name:    name,
		colName: getColName(options),
		genFunc: genFunc,
	}
}

type boolAttr struct {
	val     bool
	name    string
	colName string
	process Processor
	genFunc func() bool
}

func (attr *boolAttr) Process(procFunc Processor) Attributer {
	attr.process = procFunc
	return attr
}

func (attr boolAttr) GetVal() interface{} {
	return attr.val
}

func (attr *boolAttr) SetVal(val interface{}) error {
	realVal, ok := val.(bool)
	if !ok {
		return fmt.Errorf("set attribute val: val %+v is not bool", val)
	}

	attr.val = realVal
	return nil
}

func (attr boolAttr) ColName() string {
	return attr.colName
}

func (boolAttr) Kind() AttrType {
	return BoolAttr
}

func (attr *boolAttr) Gen(data interface{}) (interface{}, error) {
	attr.val = attr.genFunc()
	if attr.process != nil {
		if err := attr.process(attr, data); err != nil {
			return nil, err
		}
	}
	return attr.val, nil
}

func (attr boolAttr) Name() string {
	return attr.name
}

func Bytes(name string, genFunc func() []byte, options ...string) Attributer {
	return &bytesAttr{
		name:    name,
		colName: getColName(options),
		genFunc: genFunc,
	}
}

type bytesAttr struct {
	name    string
	colName string
	val     []byte
	genFunc func() []byte
	process Processor
}

func (attr *bytesAttr) Process(procFunc Processor) Attributer {
	attr.process = procFunc
	return attr
}

func (attr bytesAttr) GetVal() interface{} {
	return attr.val
}

func (attr *bytesAttr) SetVal(val interface{}) error {
	realVal, ok := val.([]byte)
	if !ok {
		return fmt.Errorf("set attribute val: val %+v is not []byte", val)
	}

	attr.val = realVal
	return nil
}

func (attr bytesAttr) ColName() string {
	return attr.colName
}

func (bytesAttr) Kind() AttrType {
	return BytesAttr
}

func (attr *bytesAttr) Gen(data interface{}) (interface{}, error) {
	attr.val = attr.genFunc()
	if attr.process != nil {
		if err := attr.process(attr, data); err != nil {
			return nil, err
		}
	}
	return attr.val, nil
}

func (attr bytesAttr) Name() string {
	return attr.name
}

func Factory(name string, factory Factorier, insertFirst bool, options ...string) Attributer {
	return &factoryAttr{
		name:        name,
		colName:     getColName(options),
		factory:     factory,
		insertFirst: insertFirst,
	}
}

type factoryAttr struct {
	factory     Factorier
	insertFirst bool
	val         interface{}
	name        string
	colName     string
	process     Processor
}

func (attr *factoryAttr) Process(procFunc Processor) Attributer {
	attr.process = procFunc
	return attr
}

func (attr factoryAttr) GetVal() interface{} {
	return attr.val
}

func (attr *factoryAttr) SetVal(val interface{}) error {
	attr.val = val
	return nil
}

func (attr factoryAttr) ColName() string {
	return attr.colName
}

func (attr factoryAttr) Name() string {
	return attr.name
}

func (factoryAttr) Kind() AttrType {
	return FactoryAttr
}

func (attr *factoryAttr) Gen(data interface{}) (interface{}, error) {
	factoryData, err := attr.factory.Build()
	if err != nil {
		return nil, err
	}
	attr.val = factoryData
	if attr.process != nil {
		if err := attr.process(attr, data); err != nil {
			return nil, err
		}
	}
	return attr.val, nil
}

func (attr factoryAttr) Insert(data interface{}) error {
	return attr.factory.Insert(data)
}

type timeAttr struct {
	begin   time.Time
	end     time.Time
	name    string
	colName string
}
