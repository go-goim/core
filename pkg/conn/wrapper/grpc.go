package wrapper

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type GrpcWrapper struct {
	ConnKey string // unique key
	*grpc.ClientConn
	context.Context
	cancel context.CancelFunc
}

func WrapGrpc(ctx context.Context, cc *grpc.ClientConn, connKey string) *GrpcWrapper {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx2, cancel := context.WithCancel(ctx)
	return &GrpcWrapper{
		ConnKey:    connKey,
		ClientConn: cc,
		Context:    ctx2,
		cancel:     cancel,
	}
}

func (w *GrpcWrapper) Key() string {
	return w.ConnKey
}

func (w *GrpcWrapper) Err() error {
	state := w.GetState()
	if state != connectivity.Idle && state != connectivity.Connecting && state != connectivity.Ready {
		w.cancel()
		return fmt.Errorf("connection invalid, cur state=%v. closing connection", state)
	}

	return w.Context.Err()
}

func (w *GrpcWrapper) Close() error {
	w.cancel()
	return w.ClientConn.Close()
}
