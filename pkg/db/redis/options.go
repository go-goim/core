package redis

import (
	"time"

	configv1 "github.com/yusank/goim/api/config/v1"
)

type Option func(*options)

type options struct {
	addr         string
	password     string
	maxConns     int
	minIdleConns int
	dialTimeout  time.Duration
	idleTimeout  time.Duration
}

func WithConfig(cfg *configv1.Redis) Option {
	return func(o *options) {
		o.addr = cfg.GetAddr()
		o.password = cfg.GetPassword()
		o.maxConns = int(cfg.GetMaxConns())
		o.minIdleConns = int(cfg.GetMinIdleConns())
		o.dialTimeout = cfg.GetDialTimeout().AsDuration()
		o.idleTimeout = cfg.GetIdleTimeout().AsDuration()
	}
}

func Addr(addr string) Option {
	return func(o *options) {
		o.addr = addr
	}
}

func Password(psw string) Option {
	return func(o *options) {
		o.password = psw
	}
}

func MaxConns(i int) Option {
	return func(o *options) {
		o.maxConns = i
	}
}

func MinIdleConns(i int) Option {
	return func(o *options) {
		o.minIdleConns = i
	}
}

func DialTimeout(d time.Duration) Option {
	return func(o *options) {
		o.dialTimeout = d
	}
}

func IdleTimeout(d time.Duration) Option {
	return func(o *options) {
		o.idleTimeout = d
	}
}
