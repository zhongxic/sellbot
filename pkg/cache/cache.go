package cache

import (
	"runtime"
	"sync/atomic"
	"time"

	"github.com/zhongxic/sellbot/pkg/container"
)

type Stringer container.Stringer

type Cache[K comparable, V any] interface {
	Set(key K, value V, expiration ...time.Duration)
	Get(key K) (value V, exist bool)
	Remove(key K)
}

type cacheItem[V any] struct {
	value  V
	expire time.Time
}

type cacheImpl[K comparable, V any] struct {
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	cm                *container.ConcurrentMap[K, cacheItem[V]]
	janitor           *janitor[K, V]
}

func (c *cacheImpl[K, V]) Set(key K, value V, expiration ...time.Duration) {
	ttl := c.defaultExpiration
	if len(expiration) > 0 && expiration[0] >= 0 {
		ttl = expiration[0]
	}
	item := cacheItem[V]{
		value:  value,
		expire: time.Now().Add(ttl),
	}
	c.cm.Put(key, item)
}

func (c *cacheImpl[K, V]) Get(key K) (value V, exist bool) {
	item, ok := c.cm.Get(key)
	if !ok || time.Now().After(item.expire) {
		return
	}
	return item.value, true
}

func (c *cacheImpl[K, V]) Remove(key K) {
	c.cm.Remove(key)
}

type cacheHolder[K comparable, V any] struct {
	*cacheImpl[K, V]
}

type janitor[K comparable, V any] struct {
	ticker *time.Ticker
	quit   chan bool
	cache  *cacheImpl[K, V]
}

func (j *janitor[K, V]) start() {
	running := atomic.Bool{}
	for {
		select {
		case <-j.ticker.C:
			if running.CompareAndSwap(false, true) {
				j.cache.cm.Range(func(key K, value cacheItem[V]) bool {
					if time.Now().After(value.expire) {
						j.cache.Remove(key)
					}
					return true
				})
				running.Store(false)
			}
		case <-j.quit:
			return
		}
	}
}

func (j *janitor[K, V]) stop() {
	j.ticker.Stop()
	j.quit <- true
}

func newJanitor[K comparable, V any](cleanupInterval time.Duration, cache *cacheImpl[K, V]) *janitor[K, V] {
	return &janitor[K, V]{
		ticker: time.NewTicker(cleanupInterval),
		quit:   make(chan bool),
		cache:  cache,
	}
}

type Options struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

func NewCache[V any](options Options) Cache[string, V] {
	c := &cacheImpl[string, V]{
		defaultExpiration: options.DefaultExpiration,
		cleanupInterval:   options.CleanupInterval,
		cm:                container.NewConcurrentMap[cacheItem[V]](),
	}
	c.janitor = newJanitor(options.CleanupInterval, c)
	h := &cacheHolder[string, V]{c}
	runtime.SetFinalizer(h, func(h *cacheHolder[string, V]) {
		h.janitor.stop()
	})
	go c.janitor.start()
	return h
}

func NewStringerCache[K Stringer, V any](options Options) Cache[K, V] {
	c := &cacheImpl[K, V]{
		defaultExpiration: options.DefaultExpiration,
		cleanupInterval:   options.CleanupInterval,
		cm:                container.NewStringerConcurrentMap[K, cacheItem[V]](),
	}
	c.janitor = newJanitor(options.CleanupInterval, c)
	h := &cacheHolder[K, V]{c}
	runtime.SetFinalizer(h, func(h *cacheHolder[K, V]) {
		h.janitor.stop()
	})
	go c.janitor.start()
	return h
}
