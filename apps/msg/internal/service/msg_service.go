package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/msg/internal/app"
	"github.com/yusank/goim/apps/msg/internal/data"
)

type MqMessageService struct{}

func (s *MqMessageService) HandleMqMessage(ctx context.Context, msg ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	// msg 实际上只有一条
	err := s.handleSingleMsg(ctx, msg[0])
	if err != nil {
		return consumer.ConsumeRetryLater, nil
	}

	return consumer.ConsumeSuccess, nil
}

func (s *MqMessageService) handleSingleMsg(ctx context.Context, msg *primitive.MessageExt) error {
	req := &messagev1.PushMessageReq{}
	if err := json.Unmarshal(msg.Body, req); err != nil {
		return err
	}

	var agentId string
	str, err := app.GetApplication().Redis.Get(ctx, data.GetUserOnlineAgentKey(req.GetToUser())).Result()
	if err != nil {
		return err
	}
	agentId = str

	reg := app.GetRegister()
	cc, err := grpc.Dial(ctx, grpc.WithDiscovery(reg),
		grpc.WithEndpoint("discovry://goim.push.service"),
		grpc.WithFilter(getFilter(agentId)))
	if err != nil {
		return err
	}

	in := &messagev1.PushMessageReq{
		FromUser:        req.GetFromUser(),
		ToUser:          req.GetToUser(),
		PushMessageType: messagev1.PushMessageType_ToUser,
		ContentType:     req.GetContentType(),
		AgentId:         agentId,
	}

	out, err := messagev1.NewPushMessagerClient(cc).PushMessage(ctx, in)
	if err != nil {
		return err
	}

	if out.GetStatus() != messagev1.PushMessageRespStatus_OK {
		return fmt.Errorf(out.GetReason())
	}

	return nil
}

func getFilter(agentId string) selector.Filter {
	return func(c context.Context, nodes []selector.Node) []selector.Node {
		var filtered = make([]selector.Node, 0)
		for i, n := range nodes {
			if n.Metadata()["agentId"] == agentId {
				filtered = append(filtered, nodes[i])
			}
		}

		return filtered
	}
}
