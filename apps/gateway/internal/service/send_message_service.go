package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	responsepb "github.com/yusank/goim/api/transport/response"
	friendpb "github.com/yusank/goim/api/user/friend/v1"
	"github.com/yusank/goim/pkg/log"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/gateway/internal/app"
	"github.com/yusank/goim/pkg/mq"
)

type SendMessageService struct {
	messagev1.UnimplementedSendMessagerServer
	friendServiceConn *ggrpc.ClientConn
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

	// check is friend
	if err := s.checkCanSendMsg(ctx, req); err != nil {
		rsp.Response = responsepb.NewBaseResponseWithMessage(responsepb.Code_RelationNotExist, err.Error())
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

func (s *SendMessageService) checkCanSendMsg(ctx context.Context, req *messagev1.SendMessageReq) error {
	if err := s.checkFriendServiceConn(ctx); err != nil {
		return err
	}

	resp, err := friendpb.NewFriendServiceClient(s.friendServiceConn).IsFriend(ctx, &friendpb.BaseFriendRequest{
		Uid:       req.GetFromUser(),
		FriendUid: req.GetToUser(),
	})
	if err != nil {
		return err
	}

	if !resp.Success() {
		return resp
	}

	return nil
}

func (s *SendMessageService) checkFriendServiceConn(ctx context.Context) error {
	if s.friendServiceConn != nil {
		switch s.friendServiceConn.GetState() {
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

	cc, err := grpc.DialInsecure(ctx,
		grpc.WithDiscovery(app.GetApplication().Register),
		grpc.WithEndpoint(fmt.Sprintf("discovery://dc1/%s", app.GetApplication().Config.SrvConfig.UserService)))
	if err != nil {
		return err
	}

	s.friendServiceConn = cc
	return nil
}

func (s *SendMessageService) Broadcast(ctx context.Context, req *messagev1.SendMessageReq) (*messagev1.SendMessageResp, error) {
	rsp := new(messagev1.SendMessageResp)
	// check req params
	if err := req.Validate(); err != nil {
		rsp.Response = responsepb.NewBaseResponseWithMessage(responsepb.Code_InvalidParams, err.Error())
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
		rsp.Response = responsepb.NewBaseResponseWithError(err)
		return rsp, nil
	}

	// todo: maybe use another topic for all broadcast messages
	rs, err := app.GetApplication().Producer.SendSync(ctx, mq.NewMessage("def_topic", b))
	if err != nil {
		rsp.Response = responsepb.NewBaseResponseWithError(err)
		return rsp, nil
	}

	log.Info(rs.String())
	rsp.MsgSeq = rs.MsgID

	return rsp, nil
}
