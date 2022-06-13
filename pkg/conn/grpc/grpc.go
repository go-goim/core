package grpc

import (
	"context"

	responsepb "github.com/go-goim/api/transport/response"
	kratosGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type ClientConn struct {
	*grpc.ClientConn
	kratosClientOpts []kratosGrpc.ClientOption
	insecure         bool

	ctx    context.Context
	cancel context.CancelFunc
}

func NewClientConn(insecure bool, opts ...kratosGrpc.ClientOption) *ClientConn {
	return &ClientConn{
		insecure:         insecure,
		kratosClientOpts: opts,
	}
}

func (c *ClientConn) Connect() error {
	var dialFunc func(ctx context.Context, co ...kratosGrpc.ClientOption) (*grpc.ClientConn, error)
	if c.insecure {
		dialFunc = kratosGrpc.Dial
	} else {
		dialFunc = kratosGrpc.DialInsecure
	}

	c.ctx, c.cancel = context.WithCancel(context.Background())

	cc, err := dialFunc(c.ctx, c.kratosClientOpts...)
	if err != nil {
		return err
	}
	c.ClientConn = cc

	return nil
}

func (c *ClientConn) Close() error {
	if c.cancel != nil {
		c.cancel()
	}

	if c.ClientConn != nil {
		return c.ClientConn.Close()
	}
	return nil
}

// conn factory

type clientConnFactory struct {
	insecure bool
	opts     []kratosGrpc.ClientOption
}

var (
	ErrConnNotReady = responsepb.Code_InternalError.BaseResponseWithMessage("conn not ready")
)

func newClientConnFactory(insecure bool, opts ...kratosGrpc.ClientOption) *clientConnFactory {
	return &clientConnFactory{
		opts: opts,
	}
}

func (c *clientConnFactory) Factory() (*ClientConn, error) {
	cc := NewClientConn(c.insecure, c.opts...)
	if err := cc.Connect(); err != nil {
		return nil, err
	}

	return cc, nil
}

func (c *clientConnFactory) Ping(cc *ClientConn) error {
	if cc.GetState() != connectivity.Ready {
		return ErrConnNotReady
	}
	return nil
}
