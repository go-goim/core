package cache

import (
	"context"
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

func NewMemoryCache() Cache {
	return &memoryCache{
		size:  defaultSize,
		items: make(map[string]*memoryCacheItem, defaultSize),
		mu:    sync.RWMutex{},
	}
}

func (m *memoryCache) Get(_ context.Context, key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, ok := m.items[key]
	if !ok {
		return nil, ErrCacheMiss
	}

	if item.expireAt.Before(time.Now()) {
		go m.Delete(context.TODO(), key) //nolint:errcheck
		return nil, ErrCacheMiss
	}

	return item.value.([]byte), nil
}

func (m *memoryCache) Set(_ context.Context, key string, value []byte, expire time.Duration) error {
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

func (m *memoryCache) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, key)
	return nil
}

func (m *memoryCache) Close(_ context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// flush all items
	m.items = make(map[string]*memoryCacheItem, m.size)
	return nil
}
