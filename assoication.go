package factory

import (
	"fmt"
	"reflect"
)

type AssociationType int8

const (
	BelongsTo AssociationType = iota + 1
	HasOneOrMany
)

type Association struct {
	factory    *Factory
	fieldName  string
	colName    string
	foreignKey string
	referField string
	num        int32
}

func (as *Association) clone() *Association {
	return &Association{
		factory:    as.factory,
		fieldName:  as.fieldName,
		colName:    as.colName,
		foreignKey: as.foreignKey,
		referField: as.referField,
		num:        as.num,
	}
}

func (as *Association) ReferField(referField string) *Association {
	as.referField = referField
	return as
}

func (as *Association) FieldName(fieldName string) *Association {
	as.fieldName = fieldName
	return as
}

func (as *Association) ColumnName(colName string) *Association {
	as.colName = colName
	return as
}

func (as *Association) ForeignField(foreignKey string) *Association {
	as.foreignKey = foreignKey
	return as
}

func (as *Association) Num(num int32) *Association {
	as.num = num
	return as
}

func (as *Association) buildFieldValue(val reflect.Value) (*fieldValue, error) {
	if as.referField != "" {
		fieldVal := getFieldValue(val.Interface(), as.referField)
		if fieldVal == nil {
			return nil, fmt.Errorf("association: has many association referField(%s) not found", as.referField)
		}
		return &fieldValue{
			val:       fieldVal,
			fieldName: as.foreignKey,
			colName:   as.colName,
		}, nil
	}
	return nil, nil
}

func (as *Association) build(val reflect.Value, insert bool, parent *Factory, asType AssociationType) ([]interface{}, error) {
	objects := make([]interface{}, as.num)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := range objects {
		var (
			object interface{}
			err    error
			pass   bool
		)

		if asType == HasOneOrMany {
			fv, err := as.buildFieldValue(val)
			if err != nil {
				return nil, err
			}
			if fv != nil {
				object, err = as.factory.buildObjectFor(insert, parent, fv)
				pass = true
			}
		}
		if !pass {
			object, err = as.factory.buildObjectFor(insert, parent)
		}

		if err != nil {
			return objects, err
		}
		objects[i] = object
	}

	if len(objects) == 1 {
		if err := as.setField(val, objects[0]); err != nil {
			return objects, err
		}
	}

	if len(objects) > 1 {
		for i := range objects {
			if err := as.setSlice(val, objects[i]); err != nil {
				return objects, err
			}
		}
	}

	return objects, nil
}

func (as *Association) setSlice(val reflect.Value, dependData interface{}) error {
	field := getFieldElem(val, as.fieldName)
	dependVal := reflect.ValueOf(dependData)
	if field.Kind() != reflect.Slice {
		return fmt.Errorf("association: field(%s) type should be slice", as.fieldName)
	}
	if !field.CanSet() {
		return fmt.Errorf("association: field(%s) is unsettable", as.fieldName)
	}

	// if element of slice is not pointer
	if field.Type().Elem().Kind() != reflect.Ptr {
		dependVal = dependVal.Elem()
	}

	newField := field
	if newField.Cap() == 0 {
		newField = reflect.MakeSlice(field.Type(), 0, 1)
	}
	newField = reflect.Append(newField, dependVal)
	field.Set(newField)
	return nil
}

func (as *Association) setField(val reflect.Value, dependData interface{}) error {
	field := getFieldElem(val, as.fieldName)
	dependVal := getElem(dependData)
	if !field.CanSet() {
		return fmt.Errorf("association: field(%s) is unsettable", as.fieldName)
	}
	field.Set(dependVal)
	return nil
}

type fieldValue struct {
	fieldName string
	colName   string
	val       interface{}
}

func (fv fieldValue) SetupObject(val reflect.Value) error {
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("association: object should be a pointer")
	}
	val = val.Elem()
	field := getFieldElem(val, fv.fieldName)
	fieldVal := getElem(fv.val)
	if !field.CanSet() {
		return fmt.Errorf("association: field(%s) is unsettable", fv.fieldName)
	}
	field.Set(fieldVal)
	return nil
}

func NewAssociations() *Associations {
	return &Associations{
		belongsTo:    make([]*Association, 0, 1),
		hasOneOrMany: make([]*Association, 0, 1),
	}
}

type Associations struct {
	belongsTo    []*Association
	hasOneOrMany []*Association
}

func (ass *Associations) clone() *Associations {
	belongsTo := make([]*Association, len(ass.belongsTo))
	copy(belongsTo, ass.belongsTo)
	hasOneOrMany := make([]*Association, len(ass.hasOneOrMany))
	copy(hasOneOrMany, ass.hasOneOrMany)
	return &Associations{
		belongsTo:    belongsTo,
		hasOneOrMany: hasOneOrMany,
	}
}

func (ass *Associations) addBelongsTo(as *Association) {
	ass.belongsTo = append(ass.belongsTo, as)
}

func (ass *Associations) addHasOneOrMany(as *Association) {
	ass.hasOneOrMany = append(ass.hasOneOrMany, as)
}

func (ass Associations) clear() {
	ass.belongsTo = ass.belongsTo[:0]
	ass.hasOneOrMany = ass.hasOneOrMany[:0]
}

func (ass Associations) buildBelongsTo(val reflect.Value, insert bool, parent *Factory) (map[string]interface{}, error) {
	columnValues := make(map[string]interface{})

	for i := range ass.belongsTo {
		as := ass.belongsTo[i]
		objects, err := as.build(val, insert, parent, BelongsTo)
		if err != nil {
			return columnValues, err
		}
		if len(objects) == 0 {
			return columnValues, fmt.Errorf("association: association object(%s) is empty", as.fieldName)
		}
		if insert && (as.foreignKey == "" || as.colName == "") {
			return columnValues, fmt.Errorf("association: insert belongTo object(%s) foreignKey or columnName is empty", as.fieldName)
		}
		if insert {
			value := getFieldValue(objects[0], as.foreignKey)
			if value == nil {
				return columnValues, fmt.Errorf("association: belongTo object(%s) foreignKey(%s) is incorrect", as.fieldName, as.foreignKey)
			}
			columnValues[as.colName] = value
		}
	}
	return columnValues, nil
}

func (ass Associations) buildHasOneOrMany(val reflect.Value, insert bool, parent *Factory) error {
	for i := range ass.hasOneOrMany {
		as := ass.hasOneOrMany[i]
		_, err := as.build(val, insert, parent, HasOneOrMany)
		if err != nil {
			return err
		}
	}
	return nil
}
