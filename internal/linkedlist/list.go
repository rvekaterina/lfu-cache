package linkedlist

import "iter"

// Node of list
// O(1) memory
type Node[V any] struct {
	// next and prev pointers in the doubly-linked list of elements
	next  *Node[V]
	prev  *Node[V]
	Value V
}

// Next returns the next node pointer after the given node pointer
func (n *Node[V]) Next() *Node[V] {
	return n.next
}

// Prev returns the prev node pointer of the given node pointer
func (n *Node[V]) Prev() *Node[V] {
	return n.prev
}

// List - O(size) memory
type List[V any] interface {
	// Move moves element node before element start
	// without memory allocations. Size not changes
	// O(1)
	Move(node, start *Node[V])

	// Push inserts a new element with the given value before node start and increments l.size.
	// Returns a new node which was added
	// O(1)
	Push(value V, start *Node[V]) *Node[V]

	// Pop deletes the last node from list and decrements l.size
	// O(1)
	Pop()

	// Remove deletes the given node from list and decrements l.size
	// O(1)
	Remove(node *Node[V])

	// Size returns number of elements in list
	// O(1)
	Size() int

	// Front returns the first node of list or nil if list is empty
	// If size == 0 function will return nil
	// O(1)
	Front() *Node[V]

	// Back returns the last node of list or nil if list is empty
	// If size == 0 function will return nil
	// O(1)
	Back() *Node[V]

	// Head returns the fictitious node of list
	// O(1)
	Head() *Node[V]

	// All returns the iterator
	// O(size)
	All() iter.Seq[V]
}

// listImpl represents a doubly linked list implementation. It is implemented as a ring.
type listImpl[V any] struct {
	head *Node[V]
	size int
}

func New[V any]() List[V] {
	l := &listImpl[V]{
		head: &Node[V]{},
	}
	l.head.prev = l.head
	l.head.next = l.head
	return l
}

func (l *listImpl[V]) Move(node, start *Node[V]) {
	node.prev.next = node.next
	node.next.prev = node.prev
	node.prev = start.prev
	node.next = start
	start.prev.next = node
	start.prev = node
}

func (l *listImpl[V]) Push(value V, start *Node[V]) *Node[V] {
	last := &Node[V]{
		next:  start,
		prev:  start.prev,
		Value: value,
	}
	last.prev.next = last
	last.next.prev = last
	l.size++
	return last
}

func (l *listImpl[V]) Pop() {
	l.Remove(l.Back())
}

func (l *listImpl[V]) Remove(node *Node[V]) {
	node.next.prev = node.prev
	node.prev.next = node.next
	node.next = nil
	node.prev = nil
	l.size--
}

func (l *listImpl[V]) Size() int {
	return l.size
}

func (l *listImpl[V]) Front() *Node[V] {
	if l.size == 0 {
		return nil
	}
	return l.head.next
}

func (l *listImpl[V]) Back() *Node[V] {
	if l.size == 0 {
		return nil
	}
	return l.head.prev
}

func (l *listImpl[V]) Head() *Node[V] {
	return l.head
}

func (l *listImpl[V]) All() iter.Seq[V] {
	return func(yield func(V) bool) {
		cur := l.Front()
		for range l.Size() {
			if !yield(cur.Value) {
				return
			}
			cur = cur.Next()
		}
	}
}
