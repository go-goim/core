package wrapper

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type GrpcWrapper struct {
	ConnKey string // unique key
	*grpc.ClientConn
}

func (w *GrpcWrapper) Key() string {
	return w.ConnKey
}

func (w *GrpcWrapper) IsClosed() bool {
	return w.GetState() == connectivity.Shutdown
}

func (w *GrpcWrapper) Reconcile() error {
	if w.IsClosed() {
		return fmt.Errorf("connection closed")
	}

	return nil
}
