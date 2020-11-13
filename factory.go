package factory

import (
	"github.com/vicxu416/gogo-factory/attr"
	"github.com/vicxu416/gogo-factory/dbutil"
)

func New(obj interface{}, attrs ...attr.Attributer) *Factory {
	fieldColumns := make(map[string]string)

	for _, a := range attrs {
		fieldColumns[a.Name()] = a.ColName()
	}

	return &Factory{
		initObj:         newConstructor(obj),
		setter:          attrs,
		fieldColumns:    fieldColumns,
		omits:           make(map[string]bool),
		insertJobsQueue: NewInsertJobQueue(),
		tempSetter:      make([]attr.Attributer, 0, 1),
		associations:    NewAssociations(),
	}
}

type Factory struct {
	table           string
	initObj         objectConstructor
	insertFunc      dbutil.InsertFunc
	setter          ObjectSetter
	tempSetter      ObjectSetter
	fieldColumns    map[string]string
	omits           map[string]bool
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
	defer f.clear()
	object, _, err := f.build(false)
	if err != nil {
		f.clear()
		panic(err)
	}
	return object
}

func (f *Factory) Build() (interface{}, error) {
	defer f.clear()
	object, _, err := f.build(false)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (f *Factory) MustInsert() interface{} {
	defer f.clear()
	object, _, err := f.build(true)
	if err != nil {
		f.clear()
		panic(err)
	}
	if err := f.insert(); err != nil {
		f.clear()
		panic(err)
	}
	return object
}

func (f *Factory) Insert() (interface{}, error) {
	defer f.clear()
	object, _, err := f.build(true)
	if err != nil {
		return nil, err
	}
	if err := f.insert(); err != nil {
		return nil, err
	}
	return object, nil
}

func (f *Factory) Omit(fields ...string) *Factory {
	for _, field := range fields {
		f.omits[field] = true
	}
	return f
}

func (f *Factory) Attrs(attrs ...attr.Attributer) *Factory {
	f.tempSetter = attrs
	return f
}

func (f *Factory) BelongsTo(fieldName string, ass *Association) *Factory {
	f.associations.addBelongsTo(ass.FieldName(fieldName).Num(1).clone())
	return f
}

func (f *Factory) HasOne(fieldName string, ass *Association) *Factory {
	f.associations.addHasOneOrMany(ass.FieldName(fieldName).Num(1).clone())
	return f
}

func (f *Factory) HasMany(fieldName string, num int32, ass *Association) *Factory {
	f.associations.addHasOneOrMany(ass.FieldName(fieldName).Num(num).clone())
	return f
}

func (f *Factory) ToAssociation() *Association {
	return &Association{
		factory: f.Clone(),
	}
}

func (f *Factory) Clone() *Factory {
	return &Factory{
		table:           f.table,
		initObj:         f.initObj,
		setter:          f.setter,
		tempSetter:      f.tempSetter,
		fieldColumns:    f.fieldColumns,
		omits:           f.omits,
		insertJobsQueue: NewInsertJobQueue(),
		associations:    f.associations.clone(),
	}
}

func (f *Factory) build(insert bool, fieldValues ...*fieldValue) (interface{}, *dbutil.InsertJob, error) {
	var (
		val       = f.initObj()
		err       error
		insertJob = &dbutil.InsertJob{}
	)

	err = f.setter.SetupObject(val, f.omits)
	if err != nil {
		return nil, nil, err
	}
	err = f.tempSetter.SetupObject(val, f.omits)
	if err != nil {
		return nil, nil, err
	}

	fieldColumns := f.getFieldColumns()
	for i := range fieldValues {
		fv := fieldValues[i]
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

	return val.Interface(), insertJob, nil
}

func (f *Factory) getFieldColumns() map[string]string {
	clonedFieldColumns := make(map[string]string)
	for k, v := range f.fieldColumns {
		clonedFieldColumns[k] = v
	}
	return clonedFieldColumns
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

func (f *Factory) clear() {
	f.insertJobsQueue.clear()
	f.associations.clear()
	f.tempSetter = make([]attr.Attributer, 0, 1)
	f.omits = make(map[string]bool)
}

func (f *Factory) getInsertFunc() dbutil.InsertFunc {
	if f.insertFunc != nil {
		return f.insertFunc
	}
	return options.InsertFunc
}

func (f *Factory) buildObjectFor(insert bool, other *Factory, fieldValues ...*fieldValue) (interface{}, error) {
	defer f.clear()
	data, _, err := f.build(insert, fieldValues...)
	if err != nil {
		return nil, err
	}
	if insert {
		other.insertJobsQueue.q.Enqueue(f.insertJobsQueue.q.head)
	}
	return data, nil
}
