package session

import (
	"fmt"
	"time"

	"github.com/zhongxic/sellbot/pkg/cache"
)

const (
	DefaultCleanupInterval = 5 * time.Minute
)

type Manager interface {
	Put(id string, session *Session)
	Get(id string) *Session
	Invalidate(id string)
}

type inMemoryManager struct {
	store cache.Cache[string, *Session]
}

func (m *inMemoryManager) Put(id string, session *Session) {
	m.store.Set(id, session)
}

func (m *inMemoryManager) Get(id string) *Session {
	value, _ := m.store.Get(id)
	return value
}

func (m *inMemoryManager) Invalidate(id string) {
	m.store.Remove(id)
}

type Options struct {
	Repository string
	Expiration time.Duration
}

func NewManager(options Options) (Manager, error) {
	switch options.Repository {
	case "memory":
		return newInMemoryManager(options)
	default:
		return nil, fmt.Errorf("invalid repository type [%v]", options.Repository)
	}
}

func newInMemoryManager(options Options) (manager Manager, err error) {
	manager = &inMemoryManager{
		store: cache.NewCache[*Session](cache.Options{
			DefaultExpiration: options.Expiration,
			CleanupInterval:   DefaultCleanupInterval,
		}),
	}
	return
}
