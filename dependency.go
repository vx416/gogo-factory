package factory

import (
	"fmt"
	"reflect"

	"github.com/vicxu416/gogo-factory/dbutil"
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

func (dep *dependency) build(data interface{}, insert bool, queue *ObjectsQueue) error {
	objects := make([]*dbutil.Object, dep.num)

	for i := range objects {
		object, err := dep.factory.build(insert)
		if err != nil {
			dep.factory.clear()
			return err
		}
		objects[i] = object
		if dep.process != nil {
			if err := dep.process(data, object.Data); err != nil {
				dep.factory.clear()
				return err
			}
		}
		queue.q.Enqueue(dep.factory.insertQueue.q.head)
		dep.factory.clear()
	}

	if len(objects) == 1 {
		if err := dep.setField(data, objects[0].Data); err != nil {
			return err
		}
	}

	if len(objects) > 1 {
		for i := range objects {
			if err := dep.setSlice(data, objects[i].Data); err != nil {
				return err
			}
		}
	}

	return nil
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
		before: &dependQueue{q: &Queue{}},
		after:  &dependQueue{q: &Queue{}},
	}
}

type DependencyManager struct {
	before *dependQueue
	after  *dependQueue
}

func (dm *DependencyManager) addBefore(depend *dependency) {
	dm.before.Enqueue(depend)
}

func (dm *DependencyManager) addAfter(depend *dependency) {
	dm.after.Enqueue(depend)
}

func (dm *DependencyManager) buildBefore(data interface{}, queue *ObjectsQueue, insert bool) error {
	scanner := dm.before.Scan()

	for depend := scanner(); depend != nil; depend = scanner() {
		err := depend.build(data, insert, queue)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dm *DependencyManager) buildAfter(data interface{}, queue *ObjectsQueue, insert bool) error {
	scanner := dm.after.Scan()

	for depend := scanner(); depend != nil; depend = scanner() {
		err := depend.build(data, insert, queue)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dm *DependencyManager) clear() {
	dm.after.clear()
	dm.before.clear()
}

type dependQueue struct {
	q *Queue
}

func (queue *dependQueue) Scan() func() *dependency {
	scan := queue.q.Scan()

	return func() *dependency {
		node := scan()
		if node == nil {
			return nil
		}

		return node.data.(*dependency)
	}

}

func (queue *dependQueue) clear() {
	curr := queue.q.head
	for curr != nil {
		next := curr.next
		depend := curr.data.(*dependency)
		if !depend.fix {
			if curr.prev == nil {
				queue.q.head = curr.next
			}
			if curr.next == nil {
				queue.q.tail = curr.prev
			}
			if curr.prev != nil {
				curr.prev.next = curr.next
			}
			if curr.next != nil {
				curr.next.prev = curr.prev
			}
			queue.q.len--
			curr.clear()
		}
		curr = next
	}
}

func (queue *dependQueue) Head() *dependency {
	if queue.q.Head() == nil {
		return nil
	}
	return queue.q.Head().data.(*dependency)
}

func (queue *dependQueue) Tail() *dependency {
	if queue.q.Tail() == nil {
		return nil
	}
	return queue.q.Tail().data.(*dependency)
}

func (queue *dependQueue) Enqueue(depend *dependency) {
	node := &Node{
		data: depend,
	}
	queue.q.Enqueue(node)
}

func (queue *dependQueue) Dequeue() *dependency {
	node := queue.q.Dequeue()
	if node == nil {
		return nil
	}
	return node.data.(*dependency)
}
