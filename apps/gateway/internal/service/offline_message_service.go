package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/gateway/internal/app"
)

type OfflineMessageService struct {
	msgServiceConn *ggrpc.ClientConn
}

var offlineMsgSrc = &OfflineMessageService{}

func GetOfflineMessageService() *OfflineMessageService {
	return offlineMsgSrc
}

func (s *OfflineMessageService) QueryOfflineMsg(ctx context.Context, req *messagev1.QueryOfflineMessageReq) (
	[]*messagev1.BriefMessage, error) {
	err := s.checkGrpcConn(ctx)
	if err != nil {
		return nil, err
	}

	rsp, err := messagev1.NewOfflineMessageClient(s.msgServiceConn).QueryOfflineMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.Response
	}

	return rsp.GetMessages(), nil
}

func (s *OfflineMessageService) checkGrpcConn(ctx context.Context) error {
	if s.msgServiceConn != nil {
		switch s.msgServiceConn.GetState() {
		case connectivity.Idle:
			return nil
		case connectivity.Connecting:
			return nil
		case connectivity.Ready:
			return nil
		default:
			// reconnect
		}
	}

	var ck = fmt.Sprintf("discovery://dc1/%s", app.GetApplication().Config.SrvConfig.MsgService)
	cc, err := grpc.DialInsecure(ctx,
		grpc.WithDiscovery(app.GetApplication().Register),
		grpc.WithEndpoint(ck),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return err
	}

	s.msgServiceConn = cc
	return nil
}
