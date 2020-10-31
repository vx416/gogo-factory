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
	head *Object
	tail *Object
	len  int
}

func (queue *ObjectsQueue) clear() {
	queue.head = nil
	queue.tail = nil
	queue.len = 0
}

func (queue *ObjectsQueue) Enqueue(object *Object) {
	queue.RPUSH(object)
}

func (queue *ObjectsQueue) Dequeue() *Object {
	return queue.LPOP()
}

func (queue *ObjectsQueue) LPUSH(object *Object) {
	if queue.head == nil {
		queue.head = object
		queue.tail = object
		queue.len++
		return
	}

	object.next = queue.head
	queue.head.prev = object
	queue.head = object
	queue.len++
}

func (queue *ObjectsQueue) RPUSH(object *Object) {
	if queue.head == nil {
		queue.head = object
		queue.tail = object
		queue.len++
		return
	}

	queue.tail.next = object
	object.prev = queue.tail
	queue.tail = object
	queue.len++
}

func (queue *ObjectsQueue) LPOP() *Object {
	if queue.head == nil {
		return nil
	}
	object := queue.head
	queue.head = object.next
	queue.len--
	return object
}

func (queue *ObjectsQueue) RPOP() *Object {
	if queue.tail == nil {
		return nil
	}
	object := queue.tail
	queue.tail = object.prev
	queue.len--
	return object
}
