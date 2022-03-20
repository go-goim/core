package service

import (
	"context"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/websocket"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/pkg/conn"
)

type PushMessager struct {
	messagev1.UnimplementedPushMessagerServer
}

func (p *PushMessager) PushMessage(ctx context.Context, req *messagev1.PushMessageReq) (resp *messagev1.PushMessageResp, err error) {
	log.Info("PUSH receive msg|", req.GetContent())
	c, ok := conn.GetConn(req.GetToUser())
	if !ok {
		log.Info("PUSH| user conn not found=", req.GetToUser())
		resp = &messagev1.PushMessageResp{
			Status: messagev1.PushMessageRespStatus_ConnectionNotFound,
			Reason: messagev1.PushMessageRespStatus_ConnectionNotFound.String(),
		}

		return
	}

	err1 := PushMessage(c.Conn, req)
	if err1 == nil {
		resp = &messagev1.PushMessageResp{Status: messagev1.PushMessageRespStatus_OK}
		return
	}

	log.Info("PUSH| push err=", err1)
	resp = &messagev1.PushMessageResp{
		Status: messagev1.PushMessageRespStatus_Unknown,
		Reason: err1.Error(),
	}

	return
}

func PushMessage(wc *websocket.Conn, message *messagev1.PushMessageReq) error {
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return wc.WriteMessage(websocket.TextMessage, b)
}
