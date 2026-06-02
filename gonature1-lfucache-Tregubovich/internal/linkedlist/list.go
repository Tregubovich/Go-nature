package linkedlist

import (
	"iter"
)

type Node[V any] struct {
	value *V
	next  *Node[V]
	prev  *Node[V]
}

func (n *Node[V]) Value() *V {
	return n.value
}

func (n *Node[V]) Next() *Node[V] {
	return n.next
}

func (n *Node[V]) Prev() *Node[V] {
	return n.prev
}

type LinkedList[V any] interface {
	Head() *Node[V]
	Tail() *Node[V]
	Size() int
	MoveAfter(node *Node[V], after *Node[V])
	PushBack(value V)
	PopBack()
	All() iter.Seq[*Node[V]]
}

type linkedListImpl[V any] struct {
	head *Node[V]
	tail *Node[V]
	size int
}

func New[V any]() *linkedListImpl[V] {
	return &linkedListImpl[V]{size: 0}
}

func (l *linkedListImpl[V]) Head() *Node[V] {
	return l.head
}

func (l *linkedListImpl[V]) Tail() *Node[V] {
	return l.tail
}

func (l *linkedListImpl[V]) Size() int {
	return l.size
}

func (l *linkedListImpl[V]) linkTo(node *Node[V], after *Node[V]) {
	node.prev = after
	node.next = after.next
	if after.next != nil {
		after.next.prev = node
	} else {
		l.head = node
	}
	after.next = node
}

func (l *linkedListImpl[V]) unlink(node *Node[V]) {
	if node.next != nil {
		node.next.prev = node.prev
		if node.prev == nil {
			l.tail = node.next
		}
	}

	if node.prev != nil {
		node.prev.next = node.next
		if node.next == nil {
			l.head = node.prev
		}
	}
}

func (l *linkedListImpl[V]) MoveAfter(node *Node[V], after *Node[V]) {
	if after == nil {
		panic("next cannot be nil")
	}
	if node == after {
		return
	}
	l.unlink(node)
	l.linkTo(node, after)
}

func (l *linkedListImpl[V]) PushBack(value V) {
	node := &Node[V]{value: &value}
	if l.tail == nil {
		l.head = node
		l.tail = node
	} else {
		l.tail.prev = node
		node.next = l.tail
		l.tail = node
	}
	l.size++
}

func (l *linkedListImpl[V]) PopBack() {
	if l.tail == nil {
		panic("list is empty")
	}
	l.unlink(l.tail)
	l.size--
}

func (l *linkedListImpl[V]) All() iter.Seq[*Node[V]] {
	return func(yield func(*Node[V]) bool) {
		for node := l.head; node != nil; node = node.prev {
			if !yield(node) {
				return
			}
		}
	}
}
