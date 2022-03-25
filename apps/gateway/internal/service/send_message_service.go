package service

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/go-kratos/kratos/v2/log"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/gateway/internal/app"
	"github.com/yusank/goim/pkg/mq"
)

type SendMessageService struct {
	messagev1.UnimplementedSendMessagerServer
}

var (
	sendMessageService     *SendMessageService
	sendMessageServiceOnce sync.Once
)

func GetSendMessageService() *SendMessageService {
	sendMessageServiceOnce.Do(func() {
		sendMessageService = new(SendMessageService)
	})

	return sendMessageService
}

func (s *SendMessageService) SendMessage(ctx context.Context, msg *messagev1.SendMessageReq) (*messagev1.SendMessageResp, error) {
	// check req params

	b, err := json.Marshal(msg)
	if err != nil {
		return err2Resp(err), err
	}

	rs, err := app.GetApplication().Producer.SendSync(ctx, mq.NewMessage("def_topic", b))
	if err != nil {
		return err2Resp(err), err
	}

	log.Info(rs.String())
	rsp := err2Resp(nil)
	rsp.MsgSeq = rs.MsgID

	return rsp, nil
}

func err2Resp(err error) *messagev1.SendMessageResp {
	if err == nil {
		return &messagev1.SendMessageResp{
			Reason: "ok",
		}
	}

	return &messagev1.SendMessageResp{
		Status: -1,
		Reason: err.Error(),
	}
}
