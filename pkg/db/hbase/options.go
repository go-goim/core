package hbase

import (
	configv1 "github.com/go-goim/api/config/v1"
)

type Option func(*option)

type option struct {
	addr string
}

func newOption() *option {
	return &option{
		addr: "",
	}
}

func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func Addr(addr string) Option {
	return func(o *option) {
		o.addr = addr
	}
}

func WithConfig(cfg *configv1.HBase) Option {
	return func(o *option) {
		o.addr = cfg.Addr
		// more options
	}
}
