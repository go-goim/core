package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	redisv8 "github.com/go-redis/redis/v8"
	ggrpc "google.golang.org/grpc"

	"github.com/go-goim/goim/pkg/consts"
	"github.com/go-goim/goim/pkg/log"

	messagev1 "github.com/go-goim/goim/api/message/v1"
	"github.com/go-goim/goim/apps/msg/internal/app"
	"github.com/go-goim/goim/pkg/conn/pool"
	"github.com/go-goim/goim/pkg/conn/wrapper"
)

type MqMessageService struct {
	rdb *redisv8.Client // remove to dao
}

var (
	mqMessageService *MqMessageService
	once             sync.Once
)

func GetMqMessageService() *MqMessageService {
	once.Do(func() {
		mqMessageService = new(MqMessageService)
		mqMessageService.rdb = app.GetApplication().Redis
	})

	return mqMessageService
}

func (s *MqMessageService) Group() string {
	return "push_msg"
}

func (s *MqMessageService) Topic() string {
	return "def_topic"
}

func (s *MqMessageService) Consume(ctx context.Context, msg ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	// msg 实际上只有一条
	err := s.handleSingleMsg(ctx, msg[0])
	if err != nil {
		log.Info("consumer error", "msg", string(msg[0].Body), "error", err)
	}

	return consumer.ConsumeSuccess, nil
}

func (s *MqMessageService) handleSingleMsg(ctx context.Context, msg *primitive.MessageExt) error {
	// PushMessageReq contains all MqMessage fields.
	req := &messagev1.MqMessage{}
	if err := json.Unmarshal(msg.Body, req); err != nil {
		return err
	}

	in := &messagev1.PushMessageReq{
		FromUser:        req.GetFromUser(),
		ToUser:          req.GetToUser(),
		PushMessageType: req.GetPushMessageType(),
		ContentType:     req.GetContentType(),
		Content:         req.GetContent(),
		MsgSeq:          msg.MsgId,
	}

	if req.GetPushMessageType() == messagev1.PushMessageType_Broadcast {
		return s.broadcast(ctx, in)
	}

	str, err := s.rdb.Get(ctx, consts.GetUserOnlineAgentKey(req.GetToUser())).Result()
	if err != nil {
		if err == redisv8.Nil {
			log.Info("user offline, put to offline queue", "user_id", req.GetToUser())
			return s.putToRedis(ctx, msg, in)
		}
		return err
	}

	in.AgentId = str
	cc, err := s.loadGrpcConn(ctx, in.AgentId)
	if err != nil {
		return err
	}

	out, err := messagev1.NewPushMessagerClient(cc).PushMessage(ctx, in)
	if err != nil {
		log.Info("MSG send msg err=", err)
		return err
	}

	if !out.Success() {
		return out
	}

	return nil
}

func (s *MqMessageService) broadcast(ctx context.Context, req *messagev1.PushMessageReq) error {
	list, err := app.GetApplication().Register.GetService(ctx, app.GetApplication().Config.SrvConfig.PushService)
	if err != nil {
		return err
	}

	for _, instance := range list {
		for _, ep := range instance.Endpoints {
			if !strings.HasPrefix(ep, "grpc://") {
				continue
			}

			if err = s.broadcastToEndpoint(ctx, req, strings.TrimPrefix(ep, "grpc://")); err != nil {
				log.Info("broadcastToEndpoint err=", err)
			}
		}
	}

	return nil
}

func (s *MqMessageService) broadcastToEndpoint(ctx context.Context, req *messagev1.PushMessageReq, ep string) error {
	cc, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(ep))
	if err != nil {
		return err
	}

	rsp, err := messagev1.NewPushMessagerClient(cc).PushMessage(ctx, req)
	if err != nil {
		return err
	}

	if !rsp.Success() {
		return rsp
	}

	return nil
}

func (s *MqMessageService) putToRedis(ctx context.Context, ext *primitive.MessageExt, req *messagev1.PushMessageReq) error {
	msgID, err := primitive.UnmarshalMsgID([]byte(ext.MsgId))
	if err != nil {
		log.Info("unmarshal ext id err=", err)
		return err
	}
	log.Info("unmarshal ext", "host", msgID.Addr, "port", msgID.Port, "offset", msgID.Offset)

	msg := &messagev1.BriefMessage{
		FromUser:    req.GetFromUser(),
		ToUser:      req.GetToUser(),
		ContentType: req.GetContentType(),
		Content:     req.GetContent(),
		MsgSeq:      ext.MsgId,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	key := consts.GetUserOfflineQueueKey(req.GetToUser())

	// add to queue
	pp := s.rdb.Pipeline()
	_ = pp.Process(ctx, s.rdb.ZAdd(ctx, key, &redisv8.Z{
		Score:  float64(msgID.Offset),
		Member: string(body),
	}))
	// set key expire
	_ = pp.Process(ctx, s.rdb.Expire(ctx, key, consts.UserOfflineQueueKeyExpire))
	// trim old messages
	_ = pp.Process(ctx, s.rdb.ZRemRangeByRank(ctx, key, 0, -int64(consts.UserOfflineQueueMemberMax+1)))

	_, err = pp.Exec(ctx)
	if err != nil {
		log.Info("Exec pipeline err=", err)
	}

	return nil
}

// todo: is there any better way to do this?
func (s *MqMessageService) loadGrpcConn(ctx context.Context, agentID string) (cc *ggrpc.ClientConn, err error) {
	var (
		ep = fmt.Sprintf("discovery://dc1/%s", app.GetApplication().Config.SrvConfig.PushService)
		ck = fmt.Sprintf("%s:%s", ep, agentID)
	)
	c := pool.Get(ck)
	if c != nil {
		wc := c.(*wrapper.GrpcWrapper)
		return wc.ClientConn, nil
	}

	cc, err = grpc.DialInsecure(ctx,
		grpc.WithDiscovery(app.GetApplication().Register),
		grpc.WithEndpoint(ep),
		grpc.WithFilter(getFilter(agentID)))

	if err != nil {
		return
	}

	pool.Add(wrapper.WrapGrpc(context.Background(), cc, ck))
	return
}

func getFilter(agentID string) selector.Filter {
	return func(c context.Context, nodes []selector.Node) []selector.Node {
		var filtered = make([]selector.Node, 0)
		for i, n := range nodes {
			log.Info("filter", n.ServiceName(), n.Address(), n.Metadata())
			if n.Metadata()["agentID"] == agentID {
				filtered = append(filtered, nodes[i])
			}
		}

		return filtered
	}
}
