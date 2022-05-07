package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"

	messagev1 "github.com/yusank/goim/api/message/v1"
	apiresp "github.com/yusank/goim/api/transport/response"
	"github.com/yusank/goim/apps/gateway/internal/app"
	"github.com/yusank/goim/pkg/conn/pool"
	"github.com/yusank/goim/pkg/conn/wrapper"
)

type OfflineMessageService struct {
}

var offlineMsgSrc = &OfflineMessageService{}

func GetOfflineMessageService() *OfflineMessageService {
	return offlineMsgSrc
}

func (s *OfflineMessageService) QueryOfflineMsg(ctx context.Context, req *messagev1.QueryOfflineMessageReq) (
	*apiresp.PbResponse, error) {
	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	return messagev1.NewOfflineMessageClient(cc).QueryOfflineMessage(ctx, req)
}

func (s *OfflineMessageService) loadConn(ctx context.Context) (*ggrpc.ClientConn, error) {
	var ck = "discovery://dc1/goim.msg.service"
	c := pool.Get(ck)
	if c != nil {
		wc := c.(*wrapper.GrpcWrapper)
		return wc.ClientConn, nil
	}

	cc, err := grpc.DialInsecure(ctx,
		grpc.WithDiscovery(app.GetApplication().Register),
		grpc.WithEndpoint(ck))
	if err != nil {
		return nil, err
	}

	pool.Add(wrapper.WrapGrpc(context.Background(), cc, ck))

	return cc, nil
}
