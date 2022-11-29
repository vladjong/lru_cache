package lrucache

import (
	"container/list"
	"sync"
)

type Node struct {
	key   string
	value interface{}
}

func NewNode(key string, value interface{}) *Node {
	return &Node{
		key:   key,
		value: value,
	}
}

type lruCache struct {
	mx       sync.RWMutex
	storage  map[string]*list.Element
	queue    *list.List
	capacity int
}

func New(capacity int) *lruCache {
	return &lruCache{
		storage:  make(map[string]*list.Element, capacity),
		queue:    list.New(),
		capacity: capacity,
	}
}

func (l *lruCache) Add(key string, value interface{}) error {
	l.mx.Lock()
	defer l.mx.Unlock()
	node := NewNode(key, value)
	if node, ok := l.storage[key]; ok {
		l.queue.MoveToFront(node)
		return nil
	}

	if l.queue.Len() == l.capacity {
		if node := l.queue.Back(); node != nil {
			l.queue.Remove(node)
			delete(l.storage, node.Value.(*Node).key)
		} else {
			return ErrQueueEmpty
		}
	}

	element := l.queue.PushFront(node)
	l.storage[key] = element
	return nil
}

func (l *lruCache) Get(key string) (interface{}, error) {
	l.mx.RLock()
	defer l.mx.RUnlock()
	value, ok := l.storage[key]
	if !ok {
		return nil, ErrNotFound
	}
	return value, nil
}

func (l *lruCache) Delete(key string) error {
	return nil
}
