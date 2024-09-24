package lfu

import (
	"errors"
	"iter"
	"lfucache/internal/linkedlist"
)

var ErrKeyNotFound = errors.New("key not found")

const DefaultCapacity = 5

// Cache
// O(capacity) memory
type Cache[K comparable, V any] interface {
	// Get returns the value of the key if the key exists in the cache,
	// otherwise, returns ErrKeyNotFound.
	//
	// O(1)
	Get(key K) (V, error)

	// Put updates the value of the key if present, or inserts the key if not already present.
	//
	// When the cache reaches its capacity, it should invalidate and remove the least frequently used key
	// before inserting a new item. For this problem, when there is a tie
	// (i.e., two or more keys with the same frequency), the least recently used key would be invalidated.
	//
	// O(1)
	Put(key K, value V)

	// All returns the iterator in descending order of frequency.
	// If two or more keys have the same frequency, the most recently used key will be listed first.
	//
	// O(capacity)
	All() iter.Seq2[K, V]

	// Size returns the cache size.
	//
	// O(1)
	Size() int

	// Capacity returns the cache capacity.
	//
	// O(1)
	Capacity() int

	// GetKeyFrequency returns the element's frequency if the key exists in the cache,
	// otherwise, returns ErrKeyNotFound.
	//
	// O(1)
	GetKeyFrequency(key K) (int, error)
}

// cacheImpl represents LFU cache implementation
// LFU cache represents blocks. In 1 block elements have the same frequency.
// The closer the element is to the beginning of the block, the least recently it has been used.
// Block is specified by its start pointer in general linked list and size (and frequency of elements).
// The structure of Cache is:
// 1. elemList - Linked list for all elements
// 2. keyToElement - map to get element by using its key
// 3. freqToStart - map to get the start pointer to the beginning of the block by using frequency of elements there
// 4. freqToCount - map to get number of elements in block by using frequency of elements there
// 5. capacity - can be set by user, otherwise it will be DefaultCapacity
// 6. defaultValue
type cacheImpl[K comparable, V any] struct {
	elemList     linkedlist.List[*element[K, V]]
	keyToElement map[K]*linkedlist.Node[*element[K, V]]
	freqToStart  map[int]*linkedlist.Node[*element[K, V]]
	freqToCount  map[int]int
	capacity     int
	defaultValue V
}

type element[K comparable, V any] struct {
	key   K
	value V
	freq  int
}

func New[K comparable, V any](capacity ...int) *cacheImpl[K, V] {
	cap := DefaultCapacity
	if len(capacity) > 0 {
		cap = capacity[0]
	}
	if cap < 0 {
		panic("invalid capacity")
	}
	return &cacheImpl[K, V]{
		elemList:     linkedlist.New[*element[K, V]](),
		keyToElement: make(map[K]*linkedlist.Node[*element[K, V]], cap),
		freqToStart:  make(map[int]*linkedlist.Node[*element[K, V]], cap),
		freqToCount:  make(map[int]int, cap),
		capacity:     cap,
	}
}

func (l *cacheImpl[K, V]) moveToFront(link *linkedlist.Node[*element[K, V]]) {
	freq := link.Value.freq
	start := l.freqToStart[freq]
	l.elemList.Move(link, start)
	l.freqToStart[freq] = link
	l.freqToCount[freq]++
}

func (l *cacheImpl[K, V]) addNewBlock(link *linkedlist.Node[*element[K, V]]) {
	freq := link.Value.freq
	l.freqToStart[freq] = link
	l.freqToCount[freq] = 1
	if next, ok := l.freqToStart[freq-1]; ok {
		l.elemList.Move(link, next)
	}
}

func (l *cacheImpl[K, V]) deleteBlock(freq int) {
	delete(l.freqToCount, freq)
	delete(l.freqToStart, freq)
}

func (l *cacheImpl[K, V]) increaseFreq(link *linkedlist.Node[*element[K, V]]) {
	freq := link.Value.freq
	l.freqToCount[freq]--

	if l.freqToCount[freq] == 0 {
		l.deleteBlock(freq)
	} else if l.freqToStart[freq] == link {
		l.freqToStart[freq] = link.Next()
	}
	link.Value.freq++

	if _, ok := l.freqToStart[link.Value.freq]; ok {
		l.moveToFront(link)
	} else {
		l.addNewBlock(link)
	}
}

func (l *cacheImpl[K, V]) Get(key K) (V, error) {
	if link, ok := l.keyToElement[key]; ok {
		l.increaseFreq(link)
		return link.Value.value, nil
	}
	return l.defaultValue, ErrKeyNotFound
}

func (l *cacheImpl[K, V]) Put(key K, value V) {
	if link, ok := l.keyToElement[key]; ok {
		l.increaseFreq(link)
		link.Value.value = value
		return
	}

	if l.elemList.Size() == l.capacity {
		last := l.elemList.Back()
		freq := last.Value.freq
		l.freqToCount[freq]--

		if l.freqToStart[freq] == last {
			l.deleteBlock(freq)
		}
		delete(l.keyToElement, last.Value.key)
		l.elemList.Pop()
	}

	if start, ok := l.freqToStart[1]; ok {
		l.keyToElement[key] = l.elemList.Push(&element[K, V]{key: key, value: value, freq: 1}, start)
		l.freqToCount[1]++
	} else {
		l.keyToElement[key] = l.elemList.Push(&element[K, V]{key: key, value: value, freq: 1}, l.elemList.Head())
		l.freqToCount[1] = 1
	}

	l.freqToStart[1] = l.keyToElement[key]
}

func (l *cacheImpl[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for elem := range l.elemList.All() {
			if !yield(elem.key, elem.value) {
				return
			}
		}
	}
}

func (l *cacheImpl[K, V]) Size() int {
	return l.elemList.Size()
}

func (l *cacheImpl[K, V]) Capacity() int {
	return l.capacity
}

func (l *cacheImpl[K, V]) GetKeyFrequency(key K) (int, error) {
	if link, ok := l.keyToElement[key]; ok {
		return link.Value.freq, nil
	}
	return 0, ErrKeyNotFound
}
