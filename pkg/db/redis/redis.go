package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewRedis(opts ...Option) (*redis.Client, error) {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	cli := redis.NewClient(&redis.Options{
		Addr:         o.addr,
		Password:     o.password,
		DialTimeout:  o.dialTimeout,
		PoolSize:     o.maxConns,
		IdleTimeout:  o.idleTimeout,
		MinIdleConns: o.minIdleConns,
	})

	// add open tracing for rdb

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return cli, nil
}
