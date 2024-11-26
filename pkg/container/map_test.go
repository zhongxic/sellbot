package container

import (
	"sync"
	"testing"
)

func TestConcurrentMap(t *testing.T) {
	m := NewConcurrentMap[int, int]()
	wg := &sync.WaitGroup{}

	N := 100
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
			m.Del(i)
			wg.Done()
		}()
	}
	wg.Wait()

	m.Range(func(k, v any) bool {
		t.Errorf("should not contains key [%v]", k)
		return false
	})
}
