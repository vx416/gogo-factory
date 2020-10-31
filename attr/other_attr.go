package attr

import (
	"fmt"
	"time"
)

// Attr create interface{} attributer with generated function
//  the return value of generated function must has the specific type
func Attr(name string, genFunc func() interface{}, options ...string) Attributer {
	return &attr{
		name:    name,
		colName: getColName(options),
		genFunc: genFunc,
	}
}

type attr struct {
	name    string
	colName string
	genFunc func() interface{}
	process Processor
	val     interface{}
}

func (attr attr) ColName() string {
	return attr.colName
}

func (attr attr) GetVal() interface{} {
	return attr.val
}

func (attr attr) SetVal(val interface{}) error {
	attr.val = val
	return nil
}

func (attr attr) Name() string {
	return attr.name
}

func (attr) Kind() Type {
	return UnknownAttr
}

func (attr *attr) Process(procFunc Processor) Attributer {
	attr.process = procFunc
	return attr
}

func (attr *attr) Gen(data interface{}) (interface{}, error) {
	attr.val = attr.genFunc()
	if attr.process != nil {
		if err := attr.process(attr, data); err != nil {
			return nil, err
		}
	}
	return attr.val, nil
}

func getColName(options []string) string {
	if len(options) > 0 {
		return getColName(options)
	}
	return ""
}

// Bytes create []byte attributer with generated function
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

func (bytesAttr) Kind() Type {
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

// Factory create factory attributer with givened factory object
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

func (factoryAttr) Kind() Type {
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
