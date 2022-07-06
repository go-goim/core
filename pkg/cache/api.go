package cache

import (
	"context"
	"time"
)

type Cache interface {
	// string

	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, expire time.Duration) error
	Delete(ctx context.Context, key string) error

	// set

	IsInSet(ctx context.Context, key string, member string) (bool, error)
	AddToSet(ctx context.Context, key string, member string) error
	DeleteFromSet(ctx context.Context, key string, member string) error

	// hashmap

	GetFromHash(ctx context.Context, key string, field string) ([]byte, error)
	SetToHash(ctx context.Context, key string, field string, value []byte) error
	DeleteFromHash(ctx context.Context, key string, field string) error

	Close(ctx context.Context) error
}
