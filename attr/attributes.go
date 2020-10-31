package attr

import (
	"database/sql"
	"fmt"
	"reflect"
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
	// FactoryAttr factory attribute
	BeforeFactoryAttr
	AfterFactoryAttr
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

func SetField(data interface{}, field reflect.Value, attr Attributer) (interface{}, error) {
	val, err := attr.Gen(data)
	if err != nil {
		return nil, err
	}

	if scanner, ok := isScanner(field); ok {
		if err := scanner.Scan(val); err != nil {
			return nil, fmt.Errorf("set scanner field, scan occurs error, %+v", err)
		}
		return val, nil
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

func isScanner(field reflect.Value) (sql.Scanner, bool) {
	var fieldRaw interface{}
	if field.CanAddr() {
		fieldRaw = field.Addr().Interface()
	} else {
		fieldRaw = field.Interface()
	}
	if scanner, ok := fieldRaw.(sql.Scanner); ok {
		return scanner, ok
	}
	return nil, false
}
