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
		go m.Delete(context.Background(), key) //nolint:errcheck
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

/*
 * set
 */

type set struct {
	items map[string]bool
}

func (s *set) has(key string) bool {
	_, ok := s.items[key]
	return ok
}

func (s *set) add(key string) {
	s.items[key] = true
}

func (m *memoryCache) IsInSet(_ context.Context, key string, member string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, ok := m.items[key]
	if !ok {
		return false, ErrCacheMiss
	}

	s := item.value.(*set)
	return s.has(member), nil
}

func (m *memoryCache) AddToSet(_ context.Context, key string, member string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.items[key]
	if !ok {
		return ErrCacheMiss
	}

	s := item.value.(*set)
	if s.has(member) {
		return nil
	}

	s.add(member)
	return nil
}

func (m *memoryCache) DeleteFromSet(_ context.Context, key string, member string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.items[key]
	if !ok {
		return ErrCacheMiss
	}

	s := item.value.(*set)
	if !s.has(member) {
		return nil
	}

	delete(s.items, member)
	return nil
}

/*
 * hashmap
 */

type hashmap struct {
	items map[string]map[string][]byte
}

func (h *hashmap) set(key string, field string, value []byte) {
	if _, ok := h.items[key]; !ok {
		h.items[key] = make(map[string][]byte)
	}
	h.items[key][field] = value
}

func (h *hashmap) get(key string, field string) ([]byte, bool) {
	if _, ok := h.items[key]; !ok {
		return nil, false
	}

	value, ok := h.items[key][field]
	return value, ok
}

func (m *memoryCache) GetFromHash(_ context.Context, key string, field string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, ok := m.items[key]
	if !ok {
		return nil, ErrCacheMiss
	}

	h, ok := item.value.(*hashmap)
	if !ok {
		return nil, ErrKeyType
	}

	value, ok := h.get(key, field)
	if !ok {
		return nil, ErrCacheMiss
	}

	return value, nil
}

func (m *memoryCache) SetToHash(_ context.Context, key string, field string, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	h := &hashmap{
		items: make(map[string]map[string][]byte),
	}
	item, ok := m.items[key]
	if ok {
		h, ok = item.value.(*hashmap)
		if !ok {
			return ErrKeyType
		}
	} else {
		m.items[key] = &memoryCacheItem{
			value: h,
		}
	}

	h.set(key, field, value)
	return nil
}

func (m *memoryCache) DeleteFromHash(_ context.Context, key string, field string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.items[key]
	if !ok {
		return ErrCacheMiss
	}

	h, ok := item.value.(*hashmap)
	if !ok {
		return ErrKeyType
	}

	_, ok = h.items[key]
	if !ok {
		return ErrCacheMiss
	}

	delete(h.items[key], field)
	return nil
}

func (m *memoryCache) Close(_ context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// flush all items
	m.items = make(map[string]*memoryCacheItem, m.size)
	return nil
}
