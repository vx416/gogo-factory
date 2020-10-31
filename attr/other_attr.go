package attr

import "reflect"

func Bool(name string, genFunc func() bool, options ...string) Attributer {
	return &boolAttr{
		name:    name,
		colName: options[0],
		genFunc: func(data interface{}) (bool, error) {
			return genFunc(), nil
		},
	}
}

func BoolWith(name string, genFunc func(data interface{}) (bool, error), options ...string) Attributer {
	return &boolAttr{
		name:    name,
		colName: options[0],
		genFunc: genFunc,
	}
}

type boolAttr struct {
	name    string
	colName string
	genFunc func(data interface{}) (bool, error)
}

func (attr boolAttr) ColName() string {
	return attr.colName
}

func (boolAttr) Kind() reflect.Kind {
	return reflect.Bool
}

func (attr boolAttr) Gen(data interface{}) (interface{}, error) {
	return attr.genFunc(data)
}

func (attr boolAttr) Name() string {
	return attr.name
}

func Bytes(name string, genFunc func() []byte, options ...string) Attributer {
	return &bytesAttr{
		name:    name,
		colName: options[0],
		genFunc: func(data interface{}) ([]byte, error) {
			return genFunc(), nil
		},
	}
}

func BytesWith(name string, genFunc func(data interface{}) ([]byte, error), options ...string) Attributer {
	return &bytesAttr{
		name:    name,
		colName: options[0],
		genFunc: genFunc,
	}
}

type bytesAttr struct {
	name    string
	colName string
	genFunc func(data interface{}) ([]byte, error)
}

func (attr bytesAttr) ColName() string {
	return attr.colName
}

func (bytesAttr) Kind() reflect.Kind {
	return reflect.Slice
}

func (attr bytesAttr) Gen(data interface{}) (interface{}, error) {
	return attr.genFunc(data)
}

func (attr bytesAttr) Name() string {
	return attr.name
}

func Factory(name string, factory Factorier, options ...string) Attributer {
	return &factoryAttr{
		name:    name,
		colName: options[0],
		factory: factory,
		genFunc: func(data interface{}, factoryData interface{}) (interface{}, error) {
			return factoryData, nil
		},
	}
}

func FactoryWith(name string, factory Factorier, genFunc func(data interface{}, factoryData interface{}) (interface{}, error), options ...string) Attributer {
	return &factoryAttr{
		name:    name,
		colName: options[0],
		factory: factory,
		genFunc: genFunc,
	}
}

type factoryAttr struct {
	factory Factorier
	name    string
	colName string
	genFunc func(data interface{}, factoryData interface{}) (interface{}, error)
}

func (attr factoryAttr) ColName() string {
	return attr.colName
}

func (attr factoryAttr) Name() string {
	return attr.name
}

func (factoryAttr) Kind() reflect.Kind {
	return reflect.Struct
}

func (attr factoryAttr) Gen(data interface{}) (interface{}, error) {
	factoryData, err := attr.factory.Build()
	if err != nil {
		return nil, err
	}
	return attr.genFunc(data, factoryData)
}

func (attr factoryAttr) Insert(data interface{}) error {
	return attr.factory.Insert(data)
}
