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

func (s *SendMessageService) SendMessage(ctx context.Context, req *messagev1.SendMessageReq) (*messagev1.SendMessageResp, error) {
	// check req params

	mm := &messagev1.MqMessage{
		FromUser:        req.GetFromUser(),
		ToUser:          req.GetToUser(),
		PushMessageType: messagev1.PushMessageType_User,
		ContentType:     req.GetContentType(),
		Content:         req.GetContent(),
	}

	b, err := json.Marshal(mm)
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

func (s *SendMessageService) Broadcast(ctx context.Context, req *messagev1.SendMessageReq) (*messagev1.SendMessageResp, error) {
	// check req params

	mm := &messagev1.MqMessage{
		FromUser:        req.GetFromUser(),
		ToUser:          req.GetToUser(),
		PushMessageType: messagev1.PushMessageType_Broadcast,
		ContentType:     req.GetContentType(),
		Content:         req.GetContent(),
	}

	b, err := json.Marshal(mm)
	if err != nil {
		return err2Resp(err), err
	}

	// todo: maybe use another topic for all broadcast messages
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
