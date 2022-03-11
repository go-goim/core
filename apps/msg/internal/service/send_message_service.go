package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/apache/rocketmq-client-go/v2/primitive"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/msg/internal/app"
)

type SendMessageService struct {
	*messagev1.UnimplementedSendMeesagerServer
}

func (s *SendMessageService) SendMessage(ctx context.Context, msg *messagev1.SendMessageReq) (*messagev1.SendMessageResp, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	pm := primitive.NewMessage("defaultTopic", b)

	rs, err := app.GetApplication().Producer.SendSync(ctx, pm)
	if err != nil {
		return nil, err
	}

	log.Println(rs.MsgID)
	return nil, nil
}
