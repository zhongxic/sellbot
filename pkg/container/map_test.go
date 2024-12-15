package container

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
)

func TestConcurrentMapPutGetRemoveRange(t *testing.T) {
	m := NewConcurrentMap[int]()
	wg := &sync.WaitGroup{}

	N := 1000
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			key := fmt.Sprintf("%v", i)
			m.Put(key, i)
			wg.Done()
		}()
	}
	wg.Wait()

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			key := fmt.Sprintf("%v", i)
			val, ok := m.Get(key)
			if !ok {
				t.Errorf("should contains key [%v]", key)
			}
			if val != i {
				t.Errorf("expected value of key [%v] is [%v] actual [%v]", key, i, val)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			key := fmt.Sprintf("%v", i)
			m.Remove(key)
			wg.Done()
		}()
	}
	wg.Wait()

	m.Range(func(key string, value int) bool {
		t.Errorf("should not contains key [%v]", key)
		return false
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
