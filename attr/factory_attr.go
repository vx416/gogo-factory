package attr

// // Factorier define factory object in factory attribute
// type Factorier interface {
// 	Build() (interface{}, error)
// 	Insert(data interface{}) error
// }

// // Factory create factory attributer with givened factory object
// func Factory(name string, factory Factorier, insertFirst bool, options ...string) Attributer {
// 	return &factoryAttr{
// 		name:        name,
// 		colName:     getColName(options),
// 		factory:     factory,
// 		insertFirst: insertFirst,
// 	}
// }

// type factoryAttr struct {
// 	factory     Factorier
// 	val         interface{}
// 	name        string
// 	colName     string
// 	process     Processor
// 	insertFirst bool
// }

// func (attr *factoryAttr) Process(procFunc Processor) Attributer {
// 	attr.process = procFunc
// 	return attr
// }

// func (attr factoryAttr) GetVal() interface{} {
// 	return attr.val
// }

// func (attr *factoryAttr) SetVal(val interface{}) error {
// 	attr.val = val
// 	return nil
// }

// func (attr factoryAttr) ColName() string {
// 	return attr.colName
// }

// func (attr factoryAttr) Name() string {
// 	return attr.name
// }

// func (attr factoryAttr) Kind() Type {
// 	if attr.insertFirst {
// 		return BeforeFactoryAttr
// 	}
// 	return AfterFactoryAttr
// }

// func (attr *factoryAttr) Gen(data interface{}) (interface{}, error) {
// 	factoryData, err := attr.factory.Build()
// 	if err != nil {
// 		return nil, err
// 	}
// 	attr.val = factoryData
// 	if attr.process != nil {
// 		if err := attr.process(attr, data); err != nil {
// 			return nil, err
// 		}
// 	}
// 	return attr.val, nil
// }

// func (attr factoryAttr) Insert(data interface{}) error {
// 	return attr.factory.Insert(data)
// }
