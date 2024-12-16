package container

import (
	"encoding/json"
	"fmt"
	"sync"
)

const defaultShardCount = 32

type Stringer interface {
	fmt.Stringer
	comparable
}

type ConcurrentMap[K comparable, V any] struct {
	shardCount int
	shards     []*mapShard[K, V]
	hash       func(k K) uint32
}

type mapShard[K comparable, V any] struct {
	sync.RWMutex
	items map[K]V
}

func NewConcurrentMap[V any](shardCount ...int) *ConcurrentMap[string, V] {
	count := determineShardCount(shardCount)
	cm := &ConcurrentMap[string, V]{
		shardCount: count,
		shards:     make([]*mapShard[string, V], count),
		hash:       fnv32,
	}
	for i := 0; i < count; i++ {
		cm.shards[i] = &mapShard[string, V]{items: make(map[string]V)}
	}
	return cm
}

func NewStringerConcurrentMap[K Stringer, V any](shardCount ...int) *ConcurrentMap[K, V] {
	count := determineShardCount(shardCount)
	cm := &ConcurrentMap[K, V]{
		shardCount: count,
		shards:     make([]*mapShard[K, V], count),
		hash:       strfnv32[K],
	}
	for i := 0; i < count; i++ {
		cm.shards[i] = &mapShard[K, V]{items: make(map[K]V)}
	}
	return cm
}

func determineShardCount(shardCount []int) int {
	count := defaultShardCount
	if len(shardCount) > 0 && shardCount[0] > 0 {
		count = shardCount[0]
	}
	return count
}

func strfnv32[K Stringer](key K) uint32 {
	return fnv32(key.String())
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	for _, c := range key {
		hash = (hash * 16777619) ^ uint32(c)
	}
	return hash
}

func (c *ConcurrentMap[K, V]) getShard(key K) *mapShard[K, V] {
	hash := c.hash(key)
	return c.shards[int(hash)%c.shardCount]
}

func (c *ConcurrentMap[K, V]) snapshot() map[K]V {
	m := make(map[K]V)
	for i := 0; i < c.shardCount; i++ {
		shard := c.shards[i]
		shard.RLock()
		for key, value := range shard.items {
			m[key] = value
		}
		shard.RUnlock()
	}
	return m
}

func (c *ConcurrentMap[K, V]) Put(key K, value V) {
	shard := c.getShard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.items[key] = value
}

func (c *ConcurrentMap[K, V]) Get(key K) (value V, ok bool) {
	shard := c.getShard(key)
	shard.RLock()
	defer shard.RUnlock()
	value, ok = shard.items[key]
	return
}

func (c *ConcurrentMap[K, V]) Remove(key K) {
	shard := c.getShard(key)
	shard.Lock()
	defer shard.Unlock()
	delete(shard.items, key)
}

func (c *ConcurrentMap[K, V]) Range(f func(key K, value V) bool) {
	snapshot := c.snapshot()
	for key, value := range snapshot {
		f(key, value)
	}
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
