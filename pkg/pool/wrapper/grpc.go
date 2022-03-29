package wrapper

import (
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
	w.GetState()
	return w.GetState() == connectivity.Shutdown
}
