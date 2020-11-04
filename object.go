package factory

import "github.com/vicxu416/gogo-factory/dbutil"

type ObjectsQueue struct {
	q *Queue
}

func (queue *ObjectsQueue) clear() {
	queue.q.clear()
}

func (queue *ObjectsQueue) Enqueue(object *dbutil.Object) {
	node := &Node{
		data: object,
	}
	queue.q.Enqueue(node)
}

func (queue *ObjectsQueue) Dequeue() *dbutil.Object {
	node := queue.q.Dequeue()
	if node == nil {
		return nil
	}
	return node.data.(*dbutil.Object)
}
