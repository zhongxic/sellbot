package container

import "sync"

type ConcurrentMap[K comparable, V any] struct {
	m sync.Map
}

func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		m: sync.Map{},
	}
}

func (c *ConcurrentMap[K, V]) Get(key K) (value V, present bool) {
	val, ok := c.m.Load(key)
	if !ok {
		var zero V
		return zero, ok
	}
	return val.(V), ok
}

func (c *ConcurrentMap[K, V]) Put(key K, value V) {
	c.m.Store(key, value)
}

func (c *ConcurrentMap[K, V]) Del(key K) {
	c.m.Delete(key)
}

func (c *ConcurrentMap[K, V]) Clear() {
	c.m.Clear()
}

func (c *ConcurrentMap[K, V]) Range(f func(key, value any) bool) {
	c.m.Range(f)
}
