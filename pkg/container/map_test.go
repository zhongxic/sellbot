package container

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
)

func TestConcurrentMap(t *testing.T) {
	m := NewConcurrentMap[int]()
	wg := &sync.WaitGroup{}

	N := 1000
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			key := fmt.Sprintf("%v", i)
			m.Put(key, i)
		}()
	}

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			key := fmt.Sprintf("%v", i)
			value, ok := m.Get(key)
			if ok && value != i {
				t.Errorf("expected value of key [%v] is [%v] actual [%v]", key, i, value)
			}
		}()
	}

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			key := fmt.Sprintf("%v", i)
			m.Remove(key)
		}()
	}

	wg.Wait()
}

func TestConcurrentMapRemoveWhenRange(t *testing.T) {
	m := NewConcurrentMap[int]()
	N := 1000
	for i := 0; i < N; i++ {
		key := fmt.Sprintf("%v", i)
		m.Put(key, i)
	}
	m.Range(func(key string, value int) bool {
		m.Remove(key)
		return true
	})
}

func TestConcurrentMapMarshalAndUnmarshal(t *testing.T) {
	m := NewConcurrentMap[int]()
	N := 1000
	for i := 0; i < N; i++ {
		key := fmt.Sprintf("%v", i)
		m.Put(key, i)
	}
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	n := NewConcurrentMap[int]()
	if err = json.Unmarshal(data, n); err != nil {
		t.Fatal(err)
	}
	n.Range(func(key string, value int) bool {
		if expected, ok := m.Get(key); !ok || value != expected {
			t.Errorf("expected value of key [%v] is [%v] actual [%v]", key, expected, value)
		}
		return true
	})
}
