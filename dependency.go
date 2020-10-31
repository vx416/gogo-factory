package factory

import (
	"fmt"
	"reflect"
)

type Processor func(data interface{}, dependency interface{}) error

type dependency struct {
	factory *Factory
	colName string
	field   string
	process Processor
	num     int
	fix     bool
}

func (dep *dependency) build(data interface{}, insert bool, queue *ObjectsQueue) ([]*Object, error) {
	objects := make([]*Object, dep.num)

	for i := range objects {
		object, err := dep.factory.build(insert)
		if err != nil {
			return nil, err
		}
		objects[i] = object
		if dep.process != nil {
			if err := dep.process(data, object.data); err != nil {
				return nil, err
			}
		}
		dep.factory.setFixFields(object.data, object.colVals)
		queue.Enqueue(dep.factory.insertQueue.head)
		dep.factory.clear()
	}

	if len(objects) == 1 {
		if err := dep.setField(data, objects[0].data); err != nil {
			return nil, err
		}
	}

	if len(objects) > 1 {
		for i := range objects {
			if err := dep.setSlice(data, objects[i].data); err != nil {
				return nil, err
			}
		}
	}

	return objects, nil
}

func (dep *dependency) setSlice(data interface{}, dependData interface{}) error {
	val := reflect.ValueOf(data)
	val = val.Elem()
	field := val.FieldByName(dep.field)
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}
	dependVal := reflect.ValueOf(dependData)
	if field.Type().Elem().Kind() != reflect.Ptr {
		dependVal = dependVal.Elem()
	}

	if field.Kind() != reflect.Slice {
		return fmt.Errorf("field(%s) type should be slice", dep.field)
	}

	newField := field
	if newField.Cap() == 0 {
		newField = reflect.MakeSlice(field.Type(), 0, 1)
	}
	newField = reflect.Append(newField, dependVal)
	field.Set(newField)
	return nil
}

func (dep *dependency) setField(data interface{}, dependData interface{}) error {
	val := reflect.ValueOf(data)
	val = val.Elem()
	field := val.FieldByName(dep.field)
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}
	dependVal := reflect.ValueOf(dependData)
	if dependVal.Kind() == reflect.Ptr {
		dependVal = dependVal.Elem()
	}

	field.Set(dependVal)
	return nil
}

func NewDepMan() *DependencyManager {
	return &DependencyManager{
		before: make([]*dependency, 0, 1),
		after:  make([]*dependency, 0, 1),
	}
}

type DependencyManager struct {
	before []*dependency
	after  []*dependency
}

func (dm *DependencyManager) addBefore(depend *dependency) {
	if cap(dm.before) == 0 {
		dm.before = make([]*dependency, 0, 1)
	}
	dm.before = append(dm.before, depend)
}

func (dm *DependencyManager) addAfter(depend *dependency) {
	if cap(dm.after) == 0 {
		dm.after = make([]*dependency, 0, 1)
	}
	dm.after = append(dm.after, depend)
}

func (dm *DependencyManager) buildBefore(data interface{}, queue *ObjectsQueue, insert bool, colVals map[string]interface{}) error {
	for i := range dm.before {
		depend := dm.before[i]
		objects, err := depend.build(data, insert, queue)
		if err != nil {
			return err
		}
		for i := range objects {
			queue.Enqueue(objects[i])
		}
		if depend.colName != "" && len(objects) == 1 {
			colVals[depend.colName] = getID(objects[0].data)
		}
	}

	return nil
}

func (dm *DependencyManager) buildAfter(data interface{}, queue *ObjectsQueue, insert bool) error {
	for i := range dm.after {
		depend := dm.after[i]
		objects, err := depend.build(data, insert, queue)
		if err != nil {
			return err
		}
		for i := range objects {
			queue.Enqueue(objects[i])
		}
	}
	return nil
}

func (dm *DependencyManager) clear() {
	newBefore := make([]*dependency, 0, len(dm.before))
	newAfter := make([]*dependency, 0, len(dm.after))
	for i := range dm.before {
		if dm.before[i].fix {
			newBefore = append(newBefore, dm.before[i])
		}
	}

	for i := range dm.after {
		if dm.after[i].fix {
			newAfter = append(newAfter, dm.after[i])
		}
	}
	dm.before = newBefore
	dm.after = newAfter
}
