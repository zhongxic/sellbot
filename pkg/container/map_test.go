package container

import (
	"encoding/json"
	"sync"
	"testing"
)

func TestConcurrentMapPutGetRemoveRange(t *testing.T) {
	m := NewConcurrentMap[int, int]()
	wg := &sync.WaitGroup{}

	N := 10000
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			m.Put(i, i)
			wg.Done()
		}()
	}
	wg.Wait()

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			val, ok := m.Get(i)
			if !ok {
				t.Errorf("should contains key [%v]", i)
			}
			if val != i {
				t.Errorf("expected value of key [%v] is [%v] actual [%v]", i, i, val)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			m.Remove(i)
			wg.Done()
		}()
	}
	wg.Wait()

	m.Range(func(key int, value int) bool {
		t.Errorf("should not contains key [%v]", key)
		return false
	})
}

func TestConcurrentMapMarshalAndUnmarshal(t *testing.T) {
	m := NewConcurrentMap[int, int]()
	m.Put(1, 1)
	m.Put(2, 2)
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	n := ConcurrentMap[int, int]{}
	if err = json.Unmarshal(data, &n); err != nil {
		t.Fatal(err)
	}
	n.Range(func(key int, value int) bool {
		if expected, ok := m.Get(key); !ok || value != expected {
			t.Errorf("expected value of key [%v] is [%v] actual [%v]", key, expected, value)
		}
		return true
	})
}
