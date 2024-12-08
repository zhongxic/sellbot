package container

import (
	"encoding/json"
	"sync"
)

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

func (c *ConcurrentMap[K, V]) Remove(key K) {
	c.m.Delete(key)
}

func (c *ConcurrentMap[K, V]) Clear() {
	c.m.Clear()
}

func (c *ConcurrentMap[K, V]) Range(f func(key K, value V) bool) {
	c.m.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

func (c *ConcurrentMap[K, V]) MarshalJSON() ([]byte, error) {
	m := make(map[K]V)
	c.Range(func(key K, value V) bool {
		m[key] = value
		return true
	})
	return json.Marshal(m)
}

func (c *ConcurrentMap[K, V]) UnmarshalJSON(data []byte) error {
	m := make(map[K]V)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	for key, value := range m {
		c.Put(key, value)
	}
	return nil
}
