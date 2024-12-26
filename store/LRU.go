package store

type Node struct {
	key   string
	value string
	prev  *Node
	next  *Node
}

type DoublyLinkedList struct {
	head *Node
	tail *Node
}

func (l *DoublyLinkedList) moveToFront(node *Node) {
	if node == l.head {
		return
	}
	l.removeNode(node)
	l.addToFront(node)
}

func (l *DoublyLinkedList) addToFront(node *Node) {
	if l.head == nil {
		l.head = node
		l.tail = node
	} else {
		l.head.prev = node
		node.next = l.head
		l.head = node
	}
}

func (l *DoublyLinkedList) removeNode(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (l *DoublyLinkedList) removeTail() *Node {
	if l.tail == nil {
		return nil
	}
	node := l.tail
	l.tail = node.prev
	l.tail.next = nil
	return node
}
