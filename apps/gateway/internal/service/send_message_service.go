package service

import (
	"context"
	"encoding/json"
	"sync"

	responsepb "github.com/yusank/goim/api/transport/response"
	"github.com/yusank/goim/pkg/log"

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
	rsp := new(messagev1.SendMessageResp)
	// check req params
	if err := req.Validate(); err != nil {
		rsp.Response = responsepb.NewBaseResponse(responsepb.Code_InvalidParams, err.Error())
		return nil, rsp.Response
	}

	mm := &messagev1.MqMessage{
		FromUser:        req.GetFromUser(),
		ToUser:          req.GetToUser(),
		PushMessageType: messagev1.PushMessageType_User,
		ContentType:     req.GetContentType(),
		Content:         req.GetContent(),
	}

	rsp, err := s.sendMessage(ctx, mm)
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.Response
	}

	return rsp, nil
}

func (s *SendMessageService) Broadcast(ctx context.Context, req *messagev1.SendMessageReq) (*messagev1.SendMessageResp, error) {
	rsp := new(messagev1.SendMessageResp)
	// check req params
	if err := req.Validate(); err != nil {
		rsp.Response = responsepb.NewBaseResponse(responsepb.Code_InvalidParams, err.Error())
		return nil, rsp.Response
	}

	mm := &messagev1.MqMessage{
		FromUser:        req.GetFromUser(),
		ToUser:          req.GetToUser(),
		PushMessageType: messagev1.PushMessageType_Broadcast,
		ContentType:     req.GetContentType(),
		Content:         req.GetContent(),
	}

	rsp, err := s.sendMessage(ctx, mm)
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.Response
	}

	return rsp, nil
}

func (s *SendMessageService) sendMessage(ctx context.Context, mm *messagev1.MqMessage) (*messagev1.SendMessageResp, error) {
	rsp := new(messagev1.SendMessageResp)

	b, err := json.Marshal(mm)
	if err != nil {
		rsp.Response = responsepb.NewBaseResponse(responsepb.Code_UnknownError, err.Error())
		return rsp, nil
	}

	// todo: maybe use another topic for all broadcast messages
	rs, err := app.GetApplication().Producer.SendSync(ctx, mq.NewMessage("def_topic", b))
	if err != nil {
		rsp.Response = responsepb.NewBaseResponse(responsepb.Code_UnknownError, err.Error())
		return rsp, nil
	}

	log.Info(rs.String())
	rsp.MsgSeq = rs.MsgID

	return rsp, nil
}
