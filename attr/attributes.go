package attr

// Type represent attributes type
type Type int8

const (
	// IntAttr int-family attribute
	IntAttr Type = iota + 1
	// UintAttr uint-family attribute
	UintAttr
	// FloatAttr float-family attribute
	FloatAttr
	// StringAttr string attribute
	StringAttr
	// BytesAttr []byte attribute
	BytesAttr
	// FactoryAttr factory attribute
	FactoryAttr
	// TimeAttr time.Time attribute
	TimeAttr
	// BoolAttr boolean attribute
	BoolAttr
	// UnknownAttr interface{} attribute
	UnknownAttr
)

// Processor define process method interface
type Processor func(attr Attributer, data interface{}) error

// Attributer define attribute interface for factory
type Attributer interface {
	Name() string
	ColName() string
	Kind() Type
	Gen(data interface{}) (interface{}, error)
	Process(process Processor) Attributer
	GetVal() interface{}
	SetVal(val interface{}) error
}

// Factorier define factory object in factory attribute
type Factorier interface {
	Build() (interface{}, error)
	Insert(data interface{}) error
}
