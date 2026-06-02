package lfucache

import (
	"errors"
	"iter"

	"github.com/igoroutine-courses/gonature.lfucache/internal/linkedlist"
)

var ErrKeyNotFound = errors.New("key not found")

const DefaultCapacity = 5

type Cache[K comparable, V any] interface {
	Get(key K) (V, error)
	Put(key K, value V)
	All() iter.Seq2[K, V]
	Size() int
	Capacity() int
	GetKeyFrequency(key K) (int, error)
}

type element[K comparable, V any] struct {
	key   K
	value V
	freq  int
}

type CacheImpl[K comparable, V any] struct {
	elements map[K]*linkedlist.Node[element[K, V]]
	lastUsed map[int]*linkedlist.Node[element[K, V]]
	list     linkedlist.LinkedList[element[K, V]]
	capacity int
}

func New[K comparable, V any](capacity ...int) *CacheImpl[K, V] {
	if len(capacity) > 0 && capacity[0] < 0 {
		panic("capacity must be greater than or equal to 0")
	}
	cache := &CacheImpl[K, V]{
		elements: make(map[K]*linkedlist.Node[element[K, V]]),
		lastUsed: make(map[int]*linkedlist.Node[element[K, V]]),
		list:     linkedlist.New[element[K, V]](),
		capacity: DefaultCapacity,
	}
	if len(capacity) > 0 {
		cache.capacity = capacity[0]
	}
	return cache
}

func (l *CacheImpl[K, V]) Get(key K) (V, error) {
	if node, ok := l.elements[key]; ok {
		l.updateFrequency(node)
		return node.Value().value, nil
	}
	return *new(V), ErrKeyNotFound
}

func (l *CacheImpl[K, V]) updateFrequency(node *linkedlist.Node[element[K, V]]) {
	freq := node.Value().freq
	l.updateOnDelete(freq, node)
	l.updateOnAdd(freq, node)
	node.Value().freq++
	l.lastUsed[freq+1] = node
}

func (l *CacheImpl[K, V]) updateOnDelete(freq int, node *linkedlist.Node[element[K, V]]) {
	if elem, ok := l.lastUsed[freq]; ok && elem == node {
		if elem.Prev() == nil || elem.Prev().Value().freq != freq {
			delete(l.lastUsed, freq)
		} else {
			l.lastUsed[freq] = elem.Prev()
		}
	}
}

func (l *CacheImpl[K, V]) updateOnAdd(freq int, node *linkedlist.Node[element[K, V]]) {
	if next, ok := l.lastUsed[freq+1]; ok {
		l.list.MoveAfter(node, next)
	} else if next, ok := l.lastUsed[freq]; ok {
		l.list.MoveAfter(node, next)
	}
}

func (l *CacheImpl[K, V]) Put(key K, value V) {
	if node, ok := l.elements[key]; ok {
		node.Value().value = value
		_, _ = l.Get(key)
		return
	}
	if len(l.elements) >= l.capacity {
		l.replaceLeastFrequency(value, key)
	} else {
		newElement := element[K, V]{
			key:   key,
			value: value,
			freq:  0,
		}
		l.list.PushBack(newElement)
	}
	l.resetElement(key)
}

func (l *CacheImpl[K, V]) resetElement(key K) {
	last := l.list.Tail()
	l.elements[key] = last
	l.lastUsed[0] = last
	_, _ = l.Get(key)
}

func (l *CacheImpl[K, V]) replaceLeastFrequency(value V, key K) {
	last := l.list.Tail()
	if l.lastUsed[last.Value().freq] == last {
		delete(l.lastUsed, last.Value().freq)
	}
	delete(l.elements, last.Value().key)
	last.Value().key = key
	last.Value().value = value
	last.Value().freq = 0
}

func (l *CacheImpl[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for node := range l.list.All() {
			if !yield(node.Value().key, node.Value().value) {
				return
			}
		}
	}
}

func (l *CacheImpl[K, V]) Size() int {
	return len(l.elements)
}

func (l *CacheImpl[K, V]) Capacity() int {
	return l.capacity
}

func (l *CacheImpl[K, V]) GetKeyFrequency(key K) (int, error) {
	if node, ok := l.elements[key]; ok {
		return node.Value().freq, nil
	}
	return 0, ErrKeyNotFound
}
