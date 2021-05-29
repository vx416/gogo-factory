package attr

import (
	"reflect"

	"github.com/vx416/gogo-factory/reflectutil"
)

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
	// TimeAttr time.Time attribute
	TimeAttr
	// BoolAttr boolean attribute
	BoolAttr
	// UnknownAttr interface{} attribute
	UnknownAttr
)

// Processor define process method interface
type Processor func(attr Attributer) error

// Attributer define attribute interface for factory
type Attributer interface {
	Name() string
	ColName() string
	Kind() Type
	Gen(data interface{}) (interface{}, error)
	Process(process Processor) Attributer
	GetVal() interface{}
	SetVal(val interface{}) error
	GetObject() interface{}
}

func SetField(data interface{}, field reflect.Value, fieldType reflect.StructField, attr Attributer) (interface{}, error) {
	val, err := attr.Gen(data)
	if err != nil {
		return nil, err
	}

	ok, err := reflectutil.TryScan(field, val)
	if ok {
		return val, err
	}

	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}

	switch attr.Kind() {
	case IntAttr:
		realVal := val.(int)
		field.SetInt(int64(realVal))
	case StringAttr:
		realVal := val.(string)
		field.SetString(realVal)
	case FloatAttr:
		realVal := val.(float64)
		field.SetFloat(realVal)
	case BoolAttr:
		realVal := val.(bool)
		field.SetBool(realVal)
	case BytesAttr:
		realVal := val.([]byte)
		field.SetBytes(realVal)
	case UintAttr:
		realVal := val.(uint)
		field.SetUint(uint64(realVal))
	case TimeAttr:
		field.Set(reflect.ValueOf(val))
	case UnknownAttr:
		field.Set(reflect.ValueOf(val))
	}

	return val, nil
}
