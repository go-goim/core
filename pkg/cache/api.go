package cache

import (
	"time"
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, expire time.Duration) error
	Delete(key string) error
	Close() error
}
