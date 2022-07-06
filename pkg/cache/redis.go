package cache

import (
	"context"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
)

// redisCache is a wrapper around the redis client that implements the
// Cache interface.
type redisCache struct {
	client *redisv8.Client
}

var _ Cache = &redisCache{}

// NewRedisCache creates a new redisCache instance.
func NewRedisCache(cli *redisv8.Client) Cache { //nolint:deadcode,unused
	return &redisCache{
		client: cli,
	}
}

func (r *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx2, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	b, err := r.client.Get(ctx2, key).Bytes()
	if err != nil {
		if err == redisv8.Nil {
			return nil, ErrCacheMiss
		}

		return nil, err
	}

	return b, nil
}

func (r *redisCache) Set(ctx context.Context, key string, value []byte, expire time.Duration) error {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx2, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	return r.client.Set(ctx2, key, value, expire).Err()
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx2, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	err := r.client.Del(ctx2, key).Err()
	if err != nil {
		if err == redisv8.Nil {
			return nil
		}
		return err
	}

	return nil
}

func (r *redisCache) IsInSet(ctx context.Context, key string, member string) (bool, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx2, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	return r.client.SIsMember(ctx2, key, member).Val(), nil
}

func (r *redisCache) AddToSet(ctx context.Context, key string, member string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx2, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	return r.client.SAdd(ctx2, key, member).Err()
}

func (r *redisCache) DeleteFromSet(ctx context.Context, key string, member string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx2, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	return r.client.SRem(ctx2, key, member).Err()
}

func (r *redisCache) GetFromHash(ctx context.Context, key string, field string) ([]byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx2, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	b, err := r.client.HGet(ctx2, key, field).Bytes()
	if err != nil {
		if err == redisv8.Nil {
			return nil, ErrCacheMiss
		}

		return nil, err
	}

	return b, nil
}

func (r *redisCache) SetToHash(ctx context.Context, key string, field string, value []byte) error {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx2, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	return r.client.HSet(ctx2, key, field, value).Err()
}

func (r *redisCache) DeleteFromHash(ctx context.Context, key string, field string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx2, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	return r.client.HDel(ctx2, key, field).Err()
}

func (r *redisCache) Close(_ context.Context) error {
	return r.client.Close()
}
