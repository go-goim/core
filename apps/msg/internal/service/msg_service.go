package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"google.golang.org/grpc/connectivity"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/msg/internal/app"
	"github.com/yusank/goim/apps/msg/internal/data"
	"github.com/yusank/goim/pkg/pool"
)

type MqMessageService struct{}

var (
	mqMessageService *MqMessageService
	once             sync.Once
)

func GetMqMessageService() *MqMessageService {
	once.Do(func() {
		mqMessageService = new(MqMessageService)
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
		log.Infof("consumer error. msg:%s,err:%v", string(msg[0].Body), err)
	}

	return consumer.ConsumeSuccess, nil
}

func (s *MqMessageService) handleSingleMsg(ctx context.Context, msg *primitive.MessageExt) error {
	req := &messagev1.PushMessageReq{}
	if err := json.Unmarshal(msg.Body, req); err != nil {
		return err
	}

	var agentID string
	str, err := app.GetApplication().Redis.Get(ctx, data.GetUserOnlineAgentKey(req.GetToUser())).Result()
	if err != nil {
		return err
	}
	agentID = str

	cc, err := s.loadGrpcConn(ctx, agentID)
	if err != nil {
		return err
	}

	in := &messagev1.PushMessageReq{
		FromUser:        req.GetFromUser(),
		ToUser:          req.GetToUser(),
		PushMessageType: messagev1.PushMessageType_User,
		ContentType:     req.GetContentType(),
		Content:         req.GetContent(),
		AgentId:         agentID,
	}

	out, err := messagev1.NewPushMessagerClient(cc).PushMessage(ctx, in)
	if err != nil {
		log.Info("MSG send msg err=", err)
		return err
	}

	if out.GetStatus() != messagev1.PushMessageRespStatus_OK {
		return fmt.Errorf(out.GetReason())
	}

	return nil
}

func (s *MqMessageService) loadGrpcConn(ctx context.Context, agentID string) (cc *ggrpc.ClientConn, err error) {
	c := pool.Get(agentID)
	if c != nil {
		wc := c.(*wrappedConn)
		return wc.ClientConn, nil
	}

	cc, err = grpc.DialInsecure(ctx,
		grpc.WithDiscovery(app.GetRegister()),
		grpc.WithEndpoint("discovery://dc1/goim.push.service"),
		grpc.WithFilter(getFilter(agentID)))
	if err != nil {
		return
	}

	pool.Add(&wrappedConn{
		ClientConn: cc,
		agentID:    agentID,
	})
	return
}

type wrappedConn struct {
	agentID string
	*ggrpc.ClientConn
}

func (w *wrappedConn) Key() string {
	return w.agentID
}

func (w *wrappedConn) IsClosed() bool {
	return w.GetState() == connectivity.Shutdown
}

func (w *wrappedConn) Reconcile() error {
	if w.IsClosed() {
		return fmt.Errorf("connection closed")
	}

	return nil
}

func getFilter(agentID string) selector.Filter {
	return func(c context.Context, nodes []selector.Node) []selector.Node {
		var filtered = make([]selector.Node, 0)
		for i, n := range nodes {
			log.Info("filter", n.ServiceName(), n.Address(), n.Metadata())
			if n.Metadata()["agentId"] == agentID {
				filtered = append(filtered, nodes[i])
			}
		}

		return filtered
	}
}
