package cache

import (
	"sync"
	"time"
)

// memoryCache is an in-memory cache.It implements the Cache interface.
// It is safe for concurrent use.
type memoryCache struct {
	size  int
	items map[string]*memoryCacheItem
	mu    sync.RWMutex
}

type memoryCacheItem struct {
	value    interface{}
	expireAt time.Time
}

var (
	_ Cache = &memoryCache{}
)

const (
	defaultSize = 1024
)

func newMemoryCache() *memoryCache {
	return &memoryCache{
		size:  defaultSize,
		items: make(map[string]*memoryCacheItem, defaultSize),
		mu:    sync.RWMutex{},
	}
}

func (m *memoryCache) Get(key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, ok := m.items[key]
	if !ok {
		return nil, ErrCacheMiss
	}

	if item.expireAt.Before(time.Now()) {
		return nil, ErrCacheMiss
	}

	return item.value.([]byte), nil
}

func (m *memoryCache) Set(key string, value []byte, expire time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.size > 0 && len(m.items) >= m.size {
		return ErrCacheFull
	}

	item := &memoryCacheItem{
		value: value,
	}

	if expire > 0 {
		item.expireAt = time.Now().Add(expire)
	}

	m.items[key] = item
	return nil
}

func (m *memoryCache) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, key)
	return nil
}

func (m *memoryCache) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// flush all items
	m.items = make(map[string]*memoryCacheItem, m.size)
	return nil
}
