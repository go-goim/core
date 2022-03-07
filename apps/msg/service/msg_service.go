package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/msg/app"
	"github.com/yusank/goim/pkg/registry"
)

type SendMessage struct {
	*messagev1.UnimplementedSendMeesagerServer
}

func (m *SendMessage) SendMessage(ctx context.Context, req *messagev1.SendMessageReq) (*messagev1.SendMessageResp, error) {
	var agentId string
	// todo: load user online status and online push server from redis
	str, err := app.GetApplication().Redis.Get(ctx, "xxx").Result()
	if err != nil {
		return nil, err
	}
	agentId = str

	reg, err := registry.GetRegisterDiscover()
	if err != nil {
		return nil, err
	}

	cc, err := grpc.Dial(ctx, grpc.WithDiscovery(reg),
		grpc.WithEndpoint("discovry://goim.push.service"),
		grpc.WithFilter(getFilter(agentId)))
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &messagev1.SendMessageResp{
		Status: int32(out.GetStatus()),
		Reason: out.GetReason(),
	}, nil
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
