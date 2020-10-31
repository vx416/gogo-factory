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
		return options[0]
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

func Time(name string, genFunc func() time.Time, options ...string) Attributer {
	return &timeAttr{
		name:    name,
		colName: getColName(options),
		genFunc: genFunc,
	}
}

type timeAttr struct {
	val     time.Time
	name    string
	colName string
	genFunc func() time.Time
	process Processor
}

func (attr *timeAttr) Process(procFunc Processor) Attributer {
	attr.process = procFunc
	return attr
}

func (attr timeAttr) GetVal() interface{} {
	return attr.val
}

func (attr *timeAttr) SetVal(val interface{}) error {
	realVal, ok := val.(time.Time)
	if !ok {
		return fmt.Errorf("set attribute val: val %+v is not time.Time", val)
	}
	attr.val = realVal
	return nil
}

func (attr timeAttr) ColName() string {
	return attr.colName
}

func (attr timeAttr) Name() string {
	return attr.name
}

func (timeAttr) Kind() Type {
	return TimeAttr
}

func (attr *timeAttr) Gen(data interface{}) (interface{}, error) {
	attr.val = attr.genFunc()
	if attr.process != nil {
		if err := attr.process(attr, data); err != nil {
			return nil, err
		}
	}
	return attr.val, nil
}

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
	genFunc func() bool
	process Processor
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

func (attr boolAttr) Name() string {
	return attr.name
}

func (boolAttr) Kind() Type {
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
