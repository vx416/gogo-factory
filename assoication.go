package gofactory

import (
	"fmt"
	"reflect"

	"github.com/vx416/gogo-factory/dbutil"

	"github.com/vx416/gogo-factory/attr"
)

type AssociationType int8

const (
	BelongsTo AssociationType = iota + 1
	HasOneOrMany
	ManyToMany
)

type joinTable struct {
	tableName string
	attrs     []attr.Attributer
}

type Association struct {
	factory         *Factory
	fieldName       string
	foreignKey      string
	foreignField    string
	associatedField string
	referField      string
	referCol        string
	joinTable       *joinTable
	num             int32
}

func (as *Association) clone() *Association {
	return &Association{
		factory:         as.factory.Clone(),
		fieldName:       as.fieldName,
		foreignKey:      as.foreignKey,
		referField:      as.referField,
		referCol:        as.referCol,
		foreignField:    as.foreignField,
		joinTable:       as.joinTable,
		associatedField: as.associatedField,
		num:             as.num,
	}
}

func (as *Association) AssociatedField(asField string) *Association {
	cloned := as.clone()
	cloned.associatedField = asField
	return cloned
}

func (as *Association) ReferColumn(referCol string) *Association {
	cloned := as.clone()
	cloned.referCol = referCol
	return cloned
}

func (as *Association) JoinTable(joinTableName string, attrs ...attr.Attributer) *Association {
	cloned := as.clone()
	cloned.joinTable = &joinTable{
		tableName: joinTableName,
		attrs:     attrs,
	}
	return cloned
}

func (as *Association) ReferField(referField string) *Association {
	cloned := as.clone()
	cloned.referField = referField
	return cloned
}

func (as *Association) FieldName(fieldName string) *Association {
	cloned := as.clone()
	cloned.fieldName = fieldName
	return cloned
}

func (as *Association) ForeignField(foreignField string) *Association {
	cloned := as.clone()
	cloned.foreignField = foreignField
	return cloned
}

func (as *Association) ForeignKey(foreignKey string) *Association {
	cloned := as.clone()
	cloned.foreignKey = foreignKey
	return cloned
}

func (as *Association) Num(num int32) *Association {
	cloned := as.clone()
	cloned.num = num
	return cloned
}

func (as *Association) buildForeignFieldValue(val reflect.Value) (*foreignFieldValue, error) {
	if as.referField != "" {
		fieldVal := getFieldValue(val.Interface(), as.referField)
		if fieldVal == nil {
			return nil, fmt.Errorf("association: has many association referField(%s) not found", as.referField)
		}
		return &foreignFieldValue{
			val:       fieldVal,
			fieldName: as.foreignField,
			colName:   as.foreignKey,
		}, nil
	}
	return nil, nil
}

func (as *Association) buildObject(insert bool, asFactory, parentFactory *Factory, foreignFV ...*foreignFieldValue) (interface{}, error) {
	data, _, err := asFactory.build(insert, foreignFV...)
	if err != nil {
		return nil, err
	}
	return data, nil
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
			fv     *foreignFieldValue
		)
		if asType == HasOneOrMany {
			fv, err = as.buildForeignFieldValue(val)
			if err != nil {
				return nil, err
			}
		}

		object, _, err = as.factory.build(insert, fv)
		if err != nil {
			return nil, err
		}

		if asType == BelongsTo {
			err := as.setForeignField(object, val)
			if err != nil {
				return nil, err
			}
		}

		if asType == ManyToMany {
			err := as.setAssociatedField(object, val)
			if err != nil {
				return nil, err
			}
		}
		objects[i] = object
	}

	if insert {
		parent.insertJobsQueue.q.Enqueue(as.factory.insertJobsQueue.q.head)
		as.factory.insertJobsQueue.clear()
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

func (as *Association) setForeignField(associatedObj interface{}, parentValue reflect.Value) error {
	var errMsg = "association(n-to-1): fields(%s), set parent's field from belongs to object, "
	if as.foreignField == "" || as.referField == "" {
		return nil
	}

	if parentValue.Kind() == reflect.Ptr {
		parentValue = parentValue.Elem()
	}
	parentField := parentValue.FieldByName(as.foreignField)
	if !parentField.IsValid() {
		return fmt.Errorf(errMsg+"parent field(%s) is invalid", as.FieldName, as.foreignField)
	}
	if !parentField.CanSet() {
		return fmt.Errorf(errMsg+"parent field(%s) is unsettable", as.FieldName, as.foreignField)
	}

	referField := getFieldElem(getElem(associatedObj), as.referField)
	if !referField.IsValid() {
		return fmt.Errorf(errMsg+"referenced field(%s) is invalid", as.FieldName, as.referField)
	}

	ok, err := attr.TryScan(parentField, referField.Interface())
	if ok {
		return err
	}

	if isPtr(parentField) && !isPtr(referField) {
		referField = referField.Addr()
	}
	if !isPtr(parentField) && isPtr(referField) {
		referField = referField.Elem()
	}
	if valuer, ok := IsValuer(referField); ok {
		value, err := valuer.Value()
		if err != nil {
			return err
		}
		referField = reflect.ValueOf(interface{}(value))
	}
	if !parentField.Type().AssignableTo(referField.Type()) {
		return fmt.Errorf(errMsg+"referenced field(%s) can't be assigned to %s", as.referField, as.foreignField)
	}
	parentField.Set(referField)
	return nil
}

func (as *Association) setAssociatedField(associatedObj interface{}, parentValue reflect.Value) error {
	errMsg := "association(m-to-m): field(%s), set parent object to associated field failed, "

	if as.associatedField == "" {
		return fmt.Errorf(errMsg+"associated field empty", as.fieldName)
	}

	associatedValue := reflect.ValueOf(associatedObj)
	associatedField := getFieldElem(associatedValue, as.associatedField)
	if !associatedField.CanSet() {
		return fmt.Errorf(errMsg+"associated field(%s) cannot set", as.fieldName, as.associatedField)
	}
	if associatedField.Kind() != reflect.Slice {
		return fmt.Errorf(errMsg+"associated field(%s) is not slice", as.fieldName, as.associatedField)
	}
	newField := associatedField
	if newField.Cap() == 0 {
		newField = reflect.MakeSlice(associatedField.Type(), 0, 1)
	}

	if newField.Type().Elem().Kind() == reflect.Ptr && parentValue.Kind() != reflect.Ptr {
		parentValue = parentValue.Addr()
	}
	if newField.Type().Elem().Kind() != reflect.Ptr && parentValue.Kind() == reflect.Ptr {
		parentValue = parentValue.Elem()
	}

	if !newField.Type().Elem().AssignableTo(parentValue.Type()) {
		return fmt.Errorf(errMsg+"parent object's type can't be assigned to associated field(%s)", as.fieldName, as.associatedField)
	}
	newField = reflect.Append(newField, parentValue)
	associatedField.Set(newField)
	return nil
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

	ok, err := attr.TryScan(field, val)
	if ok {
		return err
	}

	field.Set(dependVal)
	return nil
}

type foreignFieldValue struct {
	fieldName string
	colName   string
	val       interface{}
}

func (fv foreignFieldValue) SetupObject(val reflect.Value) error {
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("association: object should be a pointer")
	}
	val = val.Elem()
	field := getFieldElem(val, fv.fieldName)
	ok, err := attr.TryScan(field, fv.val)
	if ok {
		return err
	}

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
		manyToMany:   make([]*Association, 0, 1),
	}
}

type Associations struct {
	belongsTo    []*Association
	hasOneOrMany []*Association
	manyToMany   []*Association
}

func (ass *Associations) clone() *Associations {
	belongsTo := make([]*Association, len(ass.belongsTo))
	for i := range ass.belongsTo {
		belongsTo[i] = ass.belongsTo[i].clone()
	}
	hasOneOrMany := make([]*Association, len(ass.hasOneOrMany))
	for i := range ass.hasOneOrMany {
		hasOneOrMany[i] = ass.hasOneOrMany[i].clone()
	}
	manyToMany := make([]*Association, len(ass.manyToMany))
	for i := range ass.manyToMany {
		manyToMany[i] = ass.manyToMany[i].clone()
	}
	return &Associations{
		belongsTo:    belongsTo,
		hasOneOrMany: hasOneOrMany,
		manyToMany:   manyToMany,
	}
}

func (ass *Associations) addBelongsTo(as *Association) {
	ass.belongsTo = append(ass.belongsTo, as)
}

func (ass *Associations) addHasOneOrMany(as *Association) {
	ass.hasOneOrMany = append(ass.hasOneOrMany, as)
}

func (ass *Associations) addManyToMany(as *Association) {
	ass.manyToMany = append(ass.manyToMany, as)
}

func (ass *Associations) clear() {
	ass.belongsTo = ass.belongsTo[:0]
	ass.hasOneOrMany = ass.hasOneOrMany[:0]
	ass.manyToMany = ass.manyToMany[:0]
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
		object := objects[0]
		if insert && (as.referField == "" || as.foreignKey == "") {
			return columnValues, fmt.Errorf("association: insert belongTo object(%s) referenced field or columnName is empty", as.fieldName)
		}
		if insert {
			value := getFieldValue(object, as.referField)
			if value == nil {
				return columnValues, fmt.Errorf("association: belongTo object(%s) referField(%s) is incorrect", as.fieldName, as.referField)
			}
			columnValues[as.foreignKey] = value
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

func (ass Associations) buildManyToMany(val reflect.Value, insert bool, parent *Factory) error {
	for i := range ass.manyToMany {
		as := ass.manyToMany[i]
		objects, err := as.build(val, insert, parent, ManyToMany)
		if err != nil {
			return err
		}
		if insert {
			for _, obj := range objects {
				insertJob, err := ass.buildJoinInsertJob(as, val.Interface(), obj)
				if err != nil {
					return err
				}
				parent.insertJobsQueue.Enqueue(insertJob)
			}
		}
	}
	return nil
}

func (ass Associations) buildJoinInsertJob(as *Association, parentObj, associatedObj interface{}) (*dbutil.InsertJob, error) {
	err := ass.validateManyToManyAss(as)
	if err != nil {
		return nil, err
	}
	colValues := make(map[string]interface{})
	referVal := getFieldValue(parentObj, as.referField)
	if referVal == nil {
		return nil, fmt.Errorf("association(m-to-m): field(%s), refer field(%s) not found", as.fieldName, as.referField)
	}
	colValues[as.referCol] = referVal
	foreginVal := getFieldValue(associatedObj, as.foreignField)
	if foreginVal == nil {
		return nil, fmt.Errorf("association(m-to-m): field(%s), foreign field(%s) not found", as.fieldName, as.foreignField)
	}
	colValues[as.foreignKey] = foreginVal

	for _, a := range as.joinTable.attrs {
		val, err := a.Gen(nil)
		if err != nil {
			return nil, fmt.Errorf("association(m-to-m): field(%s), join table attribute generate value occur error, err:%+v", as.fieldName, err)
		}
		colValues[a.ColName()] = val
	}

	job := dbutil.NewJob(reflect.Value{}, colValues)
	job.SetDB(options.DB, options.Driver, as.joinTable.tableName, "")
	return job, nil
}

func (ass Associations) validateManyToManyAss(as *Association) error {
	if as.joinTable.tableName == "" {
		return fmt.Errorf("association(m-to-m): field(%s), join table is empty", as.fieldName)
	}

	if as.referField == "" || as.referCol == "" {
		return fmt.Errorf("association(m-to-m): field(%s), refer field or refer column name is empty", as.fieldName)
	}

	if as.foreignField == "" || as.foreignKey == "" {
		return fmt.Errorf("association(m-to-m): field(%s), foreign field or foreign key is empty", as.fieldName)
	}
	return nil
}
