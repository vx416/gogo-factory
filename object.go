package factory

import "fmt"

type Object struct {
	data       interface{}
	colVals    map[string]interface{}
	insert     bool
	table      string
	insertFunc InsertFunc
	next       *Object
	prev       *Object
}

func (obj *Object) Insert() error {
	if !obj.insert {
		return nil
	}

	if obj.insertFunc != nil {
		return obj.insertFunc(options.db, obj.data)
	}

	if len(obj.colVals) == 0 {
		return fmt.Errorf("insert: attributes has no column name in %s", obj.table)
	}

	return insert(options.db, obj.table, obj.colVals)
}

type ObjectsQueue struct {
	q *Queue
}

func (queue *ObjectsQueue) clear() {
	queue.q.clear()
}

func (queue *ObjectsQueue) Head() *Object {
	if queue.q.Head() == nil {
		return nil
	}
	return queue.q.Head().data.(*Object)
}

func (queue *ObjectsQueue) Tail() *Object {
	if queue.q.Tail() == nil {
		return nil
	}
	return queue.q.Tail().data.(*Object)
}

func (queue *ObjectsQueue) Enqueue(object *Object) {
	node := &Node{
		data: object,
	}
	queue.q.Enqueue(node)
}

func (queue *ObjectsQueue) Dequeue() *Object {
	node := queue.q.Dequeue()
	if node == nil {
		return nil
	}
	return node.data.(*Object)
}
