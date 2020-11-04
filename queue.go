package factory

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
