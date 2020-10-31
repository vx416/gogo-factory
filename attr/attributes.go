package attr

type AttrType int8

const (
	IntAttr AttrType = iota + 1
	UintAttr
	FloatAttr
	StringAttr
	BytesAttr
	FactoryAttr
	TimeAttr
	BoolAttr
	UnknownAttr
)

type Processor func(attr Attributer, data interface{}) error

type Attributer interface {
	Name() string
	ColName() string
	Kind() AttrType
	Gen(data interface{}) (interface{}, error)
	Process(process Processor) Attributer
	GetVal() interface{}
	SetVal(val interface{}) error
}

type Factorier interface {
	Build() (interface{}, error)
	Insert(data interface{}) error
}

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

func (attr) Kind() AttrType {
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
