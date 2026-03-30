package memory

import (
	"errors"
	"sync"
	"time"

	"opita-sync-framework/internal/cache/store"
)

type CacheStore struct {
	mu    sync.RWMutex
	items map[string]store.Item
}

func NewCacheStore() *CacheStore {
	return &CacheStore{items: map[string]store.Item{}}
}

func (s *CacheStore) Set(item store.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now().UTC()
	}
	s.items[item.Key] = item
	return nil
}

func (s *CacheStore) Get(key string) (store.Item, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, found := s.items[key]
	if !found {
		return store.Item{}, false, nil
	}
	if item.TTL > 0 && time.Since(item.CreatedAt) > item.TTL {
		return store.Item{}, false, nil
	}
	return item, true, nil
}

func (s *CacheStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, found := s.items[key]; !found {
		return errors.New("cache item not found")
	}
	delete(s.items, key)
	return nil
}

var _ store.Service = (*CacheStore)(nil)
