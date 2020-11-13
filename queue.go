package factory

import "github.com/vicxu416/gogo-factory/dbutil"

type Node struct {
	next *Node
	prev *Node
	data interface{}
}

func (node *Node) clear() {
	node.data = nil
	node.next = nil
	node.prev = nil
}

type Queue struct {
	head *Node
	tail *Node
	len  int
}

func (queue *Queue) clear() {
	queue.head = nil
	queue.tail = nil
	queue.len = 0
}

func (queue *Queue) Head() *Node {
	return queue.head
}

func (queue *Queue) Tail() *Node {
	return queue.tail
}

func (queue *Queue) Enqueue(node *Node) {
	if queue.head == nil {
		queue.head = node
		queue.tail = node
		queue.len++
		return
	}

	queue.tail.next = node
	node.prev = queue.tail
	queue.tail = node
	queue.len++
	for queue.tail.next != nil {
		queue.tail = queue.tail.next
		queue.len++
	}
}

func (queue *Queue) Dequeue() *Node {
	if queue.head == nil {
		return nil
	}
	node := queue.head
	queue.head = node.next
	queue.len--
	return node
}

func (queue *Queue) Scan() func() *Node {
	curr := queue.Head()
	return func() *Node {
		if curr == nil {
			return curr
		}
		node := curr
		curr = curr.next
		return node
	}
}

func NewInsertJobQueue() *InsertJobsQueue {
	return &InsertJobsQueue{
		q: &Queue{},
	}
}

type InsertJobsQueue struct {
	q *Queue
}

func (queue *InsertJobsQueue) clear() {
	queue.q.clear()
}

func (queue *InsertJobsQueue) Enqueue(job *dbutil.InsertJob) {
	node := &Node{
		data: job,
	}
	queue.q.Enqueue(node)
}

func (queue *InsertJobsQueue) Dequeue() *dbutil.InsertJob {
	node := queue.q.Dequeue()
	if node == nil {
		return nil
	}
	return node.data.(*dbutil.InsertJob)
}
