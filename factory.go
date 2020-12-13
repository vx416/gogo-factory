package gofactory

import (
	"fmt"
	"reflect"

	"github.com/vx416/gogo-factory/attr"
	"github.com/vx416/gogo-factory/dbutil"
)

// New construct a factory object
func New(obj interface{}, attrs ...attr.Attributer) *Factory {
	objectSetter := ObjectSetter(attrs)
	fieldColumns := objectSetter.buildFieldColumns(obj)

	return &Factory{
		initObj:         newConstructor(obj),
		setter:          objectSetter,
		fieldColumns:    fieldColumns,
		omits:           make(map[string]bool),
		only:            make(map[string]bool),
		insertJobsQueue: NewInsertJobQueue(),
		associations:    NewAssociations(),
	}
}

type Factory struct {
	table           string
	initObj         objectConstructor
	insertFunc      dbutil.InsertFunc
	setter          ObjectSetter
	fieldColumns    map[string]string
	omits           map[string]bool
	only            map[string]bool
	insertJobsQueue *InsertJobsQueue
	associations    *Associations
}

func (f *Factory) Table(tableName string) *Factory {
	f.table = tableName
	return f
}

func (f *Factory) InsertFunc(fn dbutil.InsertFunc) *Factory {
	f.insertFunc = fn
	return f
}

func (f *Factory) MustBuild() interface{} {
	object, _, err := f.build(false)
	if err != nil {
		panic(err)
	}
	return object
}

func (f *Factory) Build() (interface{}, error) {
	object, _, err := f.build(false)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (f *Factory) MustInsert() interface{} {
	object, _, err := f.build(true)
	if err != nil {
		panic(err)
	}
	if err := f.insert(); err != nil {
		panic(err)
	}
	return object
}

func (f *Factory) Insert() (interface{}, error) {
	object, _, err := f.build(true)
	if err != nil {
		return nil, err
	}
	if err := f.insert(); err != nil {
		return nil, err
	}
	return object, nil
}

func (f *Factory) MustInsertN(n int) interface{} {
	object, err := f.buildN(n, true)
	if err != nil {
		panic(err)
	}
	err = f.insert()
	if err != nil {
		panic(err)
	}
	return object
}

func (f *Factory) InsertN(n int) (interface{}, error) {
	object, err := f.buildN(n, true)
	if err != nil {
		return nil, err
	}
	err = f.insert()
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (f *Factory) MustBuildN(n int) interface{} {
	objects, err := f.buildN(n, false)
	if err != nil {
		panic(err)
	}
	return objects
}

func (f *Factory) BuildN(n int) (interface{}, error) {
	return f.buildN(n, false)
}

func (f *Factory) Omit(fields ...string) *Factory {
	cloned := f.Clone()
	for _, field := range fields {
		cloned.omits[field] = true
	}
	return cloned
}

func (f *Factory) ClearOmit() *Factory {
	cloned := f.Clone()
	cloned.omits = make(map[string]bool)
	return cloned
}

func (f *Factory) Only(fields ...string) *Factory {
	cloned := f.Clone()
	for _, field := range fields {
		cloned.only[field] = true
	}
	return cloned
}

// Attrs replace object Attributer and return the new factory
func (f *Factory) Attrs(attrs ...attr.Attributer) *Factory {
	cloned := f.Clone()
	oldAttrsMap := make(map[string]int)
	for i := range cloned.setter {
		oldAttrsMap[cloned.setter[i].Name()] = i
	}

	for i := range attrs {
		cloned.fieldColumns[attrs[i].Name()] = attrs[i].ColName()
		oldIndex := oldAttrsMap[attrs[i].Name()]
		if oldIndex != 0 {
			cloned.setter[oldIndex] = attrs[i]
		} else {
			cloned.setter = append(cloned.setter, attrs[i])
		}
	}
	return cloned
}

func (f *Factory) BelongsTo(fieldName string, ass *Association) *Factory {
	cloned := f.Clone()
	cloned.associations.addBelongsTo(ass.FieldName(fieldName).Num(1))
	return cloned
}

func (f *Factory) HasOne(fieldName string, ass *Association) *Factory {
	cloned := f.Clone()
	cloned.associations.addHasOneOrMany(ass.FieldName(fieldName).Num(1))
	return cloned
}

func (f *Factory) HasMany(fieldName string, ass *Association, num int32) *Factory {
	cloned := f.Clone()
	cloned.associations.addHasOneOrMany(ass.FieldName(fieldName).Num(num))
	return cloned
}

func (f *Factory) ManyToMany(fieldName string, ass *Association, num int32) *Factory {
	cloned := f.Clone()
	cloned.associations.addManyToMany(ass.FieldName(fieldName).Num(num))
	return cloned
}

func (f *Factory) ToAssociation() *Association {
	return &Association{
		factory: f.Clone(),
	}
}

// Clone clone a factory object
func (f *Factory) Clone() *Factory {
	clonedOmits := make(map[string]bool)
	for k, v := range f.omits {
		clonedOmits[k] = v
	}
	clonedFieldColumns := make(map[string]string)
	for k, v := range f.fieldColumns {
		clonedFieldColumns[k] = v
	}
	clonedOnly := make(map[string]bool)
	for k, v := range f.only {
		clonedOnly[k] = v
	}

	return &Factory{
		table:           f.table,
		initObj:         f.initObj,
		setter:          f.setter.clone(),
		fieldColumns:    clonedFieldColumns,
		omits:           clonedOmits,
		only:            clonedOnly,
		insertJobsQueue: NewInsertJobQueue(),
		associations:    f.associations.clone(),
	}
}

func (f *Factory) buildN(n int, insert bool) (interface{}, error) {
	if n == 0 {
		return nil, fmt.Errorf("buildN: size(n) cannot be zero")
	}
	values := make([]reflect.Value, 0, n)
	for i := 0; i < n; i++ {
		cloned := f.Clone()
		object, _, err := cloned.build(insert)
		if err != nil {
			return nil, err
		}

		if insert {
			f.insertJobsQueue.q.Enqueue(cloned.insertJobsQueue.q.head)
		}

		values = append(values, reflect.ValueOf(object))
	}

	sliceVal := makeSlice(values[0].Interface(), n)
	sliceVal = reflect.Append(sliceVal, values...)
	return sliceVal.Interface(), nil
}

func (f *Factory) build(insert bool, foreignFV ...*foreignFieldValue) (interface{}, *dbutil.InsertJob, error) {
	var (
		val       = f.initObj()
		err       error
		insertJob = &dbutil.InsertJob{}
	)

	err = f.setter.SetupObject(val, f.omits, f.only)
	if err != nil {
		return nil, nil, err
	}

	fieldColumns := f.fieldColumns
	if Opt().TagProcess != nil {
		fieldColumns = getObjectColumnNames(val, Opt().TagProcess)
	}

	for i := range foreignFV {
		fv := foreignFV[i]
		if fv == nil {
			continue
		}
		err := fv.SetupObject(val)
		if err != nil {
			return nil, nil, err
		}
		if insert {
			fieldColumns[fv.fieldName] = fv.colName
		}
	}

	belongToValues, err := f.associations.buildBelongsTo(val, insert, f)
	if err != nil {
		return nil, nil, err
	}

	if insert {
		colValues := getColumnValues(val, fieldColumns)
		for k, v := range belongToValues {
			colValues[k] = v
		}
		insertJob = dbutil.NewJob(val, colValues)
		insertJob.SetDB(options.DB, options.Driver, f.table, "")
		insertJob.SetInsertFunc(f.getInsertFunc())
		f.insertJobsQueue.Enqueue(insertJob)
	}

	err = f.associations.buildHasOneOrMany(val, insert, f)
	if err != nil {
		return nil, nil, err
	}

	err = f.associations.buildManyToMany(val, insert, f)
	if err != nil {
		return nil, nil, err
	}

	return val.Interface(), insertJob, nil
}

func (f *Factory) insert() error {
	object := f.insertJobsQueue.Dequeue()
	var err error
	for object != nil {
		err = object.Insert()
		object = f.insertJobsQueue.Dequeue()
	}
	if err != nil {
		return err
	}
	return nil
}

func (f *Factory) getInsertFunc() dbutil.InsertFunc {
	if f.insertFunc != nil {
		return f.insertFunc
	}
	return options.InsertFunc
}
