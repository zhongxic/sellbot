package cache

import (
	"runtime"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	options := Options{
		DefaultExpiration: 2000 * time.Millisecond,
		CleanupInterval:   1000 * time.Millisecond,
	}
	cache := NewCache[string](options)

	key := "key"
	value := "value"

	cache.Set(key, value)
	retrieve, exist := cache.Get(key)
	if !exist {
		t.Errorf("key [%v] should contains in cache", key)
	}
	if retrieve != value {
		t.Errorf("value not equal expected [%v] actual [%v]", value, retrieve)
	}
	time.Sleep(options.DefaultExpiration)
	_, exist = cache.Get(key)
	if exist {
		t.Errorf("key [%v] should expire after [%v] mills", key, options.DefaultExpiration.Milliseconds())
	}

	cache = nil
	runtime.GC()
}
