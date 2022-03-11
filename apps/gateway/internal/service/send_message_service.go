package service

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/transport/grpc"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/gateway/internal/app"
)

type SendMessageService struct {
}

func (s *SendMessageService) SendMessage(ctx context.Context, msg *messagev1.SendMessageReq) error {
	cc, err := grpc.Dial(ctx, grpc.WithDiscovery(app.GetRegister()),
		grpc.WithEndpoint("discovry://goim.msg.service"))
	if err != nil {
		return err
	}

	out, err := messagev1.NewSendMeesagerClient(cc).SendMessage(ctx, msg)
	if err != nil {
		return err
	}

	if out.GetStatus() != int32(messagev1.PushMessageRespStatus_OK) {
		return fmt.Errorf(out.GetReason())
	}

	return nil
}
