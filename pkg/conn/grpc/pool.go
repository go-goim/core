package grpc

import (
	"sync/atomic"

	kratosGrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/go-goim/core/pkg/errors"
)

type ConnPool struct {
	factory *clientConnFactory
	opts    *poolOptions

	lastIndex int32 // atomic
	conns     []*ClientConn
}

type poolOptions struct {
	dialInsecure bool
	dialOpts     []kratosGrpc.ClientOption
	size         int
}

func WithClientOption(opts ...kratosGrpc.ClientOption) PoolOption {
	return func(o *poolOptions) {
		o.dialOpts = append(o.dialOpts, opts...)
	}
}

func WithInsecure() PoolOption {
	return func(o *poolOptions) {
		o.dialInsecure = true
	}
}

func WithPoolSize(size int) PoolOption {
	return func(o *poolOptions) {
		o.size = size
	}
}

func newPoolOptions(opts ...PoolOption) *poolOptions {
	o := &poolOptions{
		size: 1,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

type PoolOption func(opts *poolOptions)

func NewConnPool(opts ...PoolOption) (*ConnPool, error) {
	p := &ConnPool{
		opts: newPoolOptions(opts...),
	}
	p.conns = make([]*ClientConn, p.opts.size)

	p.factory = newClientConnFactory(p.opts.dialInsecure, p.opts.dialOpts...)
	for i := 0; i < p.opts.size; i++ {
		cc, err := p.factory.Factory()
		if err != nil {
			return nil, err
		}

		p.conns[i] = cc
	}

	return p, nil
}

func (c *ConnPool) Get() (*ClientConn, error) {
	cc := c.conns[atomic.AddInt32(&c.lastIndex, 1)%int32(c.opts.size)]
	atomic.AddInt32(&c.lastIndex, 1)

	return cc, nil
}

func (c *ConnPool) Close(conn *ClientConn) error {
	return c.factory.Close(conn)
}

func (c *ConnPool) Release() error {
	var es errors.ErrorSet
	for _, conn := range c.conns {
		if err := c.factory.Close(conn); err != nil {
			es = append(es, err)
		}
	}

	c.conns = nil

	return es.Err()
}

func (c *ConnPool) Len() int {
	return len(c.conns)
}
