package session

import (
	"time"

	"github.com/zhongxic/sellbot/pkg/cache"
)

const (
	sessionCleanupInterval = 5 * time.Minute
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
	Expiration time.Duration
}

func NewInMemoryManager(options Options) Manager {
	return &inMemoryManager{
		store: cache.NewCache[*Session](cache.Options{
			DefaultExpiration: options.Expiration,
			CleanupInterval:   sessionCleanupInterval,
		}),
	}
}
