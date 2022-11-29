package cache

import (
	"container/list"
	"sync"
	"time"
)

type Node struct {
	Key   interface{}
	Value interface{}
}

func NewNode(key, value interface{}) *Node {
	return &Node{
		Key:   key,
		Value: value,
	}
}

type lruCache struct {
	storage  map[interface{}]*list.Element
	queue    *list.List
	capacity int
	size     int
	mx       sync.RWMutex
}

func New(size int) *lruCache {
	return &lruCache{
		storage:  make(map[interface{}]*list.Element, size),
		queue:    list.New(),
		capacity: 0,
		size:     size,
	}
}

func (l *lruCache) Add(key, value interface{}) error {
	l.mx.Lock()
	defer l.mx.Unlock()
	node := NewNode(key, value)
	if node, ok := l.storage[key]; ok {
		node.Value.(*Node).Value = value
		l.queue.MoveToFront(node)
		l.storage[key] = node
		return nil
	}
	if l.queue.Len() == l.size {
		if node := l.queue.Back(); node != nil {
			l.queue.Remove(node)
			delete(l.storage, node.Value.(*Node).Key)
		} else {
			return ErrQueueEmpty
		}
	}
	l.capacity += 1
	element := l.queue.PushFront(node)
	l.storage[key] = element
	return nil
}

func (l *lruCache) Get(key interface{}) (interface{}, error) {
	l.mx.RLock()
	defer l.mx.RUnlock()
	value, ok := l.storage[key]
	if !ok {
		return nil, ErrNotFound
	}
	l.queue.MoveToFront(value)
	return value.Value.(*Node).Value, nil
}

func (l *lruCache) Remove(key interface{}) error {
	l.mx.Lock()
	defer l.mx.Unlock()
	value, ok := l.storage[key]
	if !ok {
		return ErrNotFound
	}
	delete(l.storage, key)
	l.queue.Remove(value)
	l.capacity -= 1
	return nil
}

func (l *lruCache) Cap() int {
	l.mx.RLock()
	defer l.mx.RUnlock()
	return l.capacity
}

func (l *lruCache) Clear() {
	l.mx.Lock()
	defer l.mx.Unlock()
	for k := range l.storage {
		delete(l.storage, k)
	}
	l.queue.Init()
	l.capacity = 0
}

func (l *lruCache) AddWithTTL(key, value interface{}, ttl time.Duration) {
	l.Add(key, value)
	go func() {
		time.Sleep(ttl)
		l.Remove(key)
	}()

}
