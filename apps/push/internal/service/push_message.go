package service

import (
	"context"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/websocket"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/pkg/pool"
	"github.com/yusank/goim/pkg/pool/wrapper"
)

type PushMessager struct {
	messagev1.UnimplementedPushMessagerServer
}

func (p *PushMessager) PushMessage(ctx context.Context, req *messagev1.PushMessageReq) (resp *messagev1.PushMessageResp, err error) {
	log.Info("PUSH receive msg|", req.String())
	if req.GetPushMessageType() == messagev1.PushMessageType_Broadcast {
		go p.handleBroadcastAsync(ctx, req)
		resp = &messagev1.PushMessageResp{Status: messagev1.PushMessageRespStatus_OK}
		return
	}
	c := pool.Get(req.GetToUser())
	if c == nil {
		log.Info("PUSH| user conn not found=", req.GetToUser())
		resp = &messagev1.PushMessageResp{
			Status: messagev1.PushMessageRespStatus_ConnectionNotFound,
			Reason: messagev1.PushMessageRespStatus_ConnectionNotFound.String(),
		}

		return
	}

	err1 := PushMessage(c.(*wrapper.WebsocketWrapper), req)
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

func (p *PushMessager) handleBroadcastAsync(ctx context.Context, req *messagev1.PushMessageReq) {
	_ = pool.Range(func(c pool.Conn) error {
		// todo use queued worker
		go func() {
			if err := PushMessage(c.(*wrapper.WebsocketWrapper), req); err != nil {
				log.Info("PushMessage err=", err)
			}
		}()

		return nil
	})
}

func PushMessage(ww *wrapper.WebsocketWrapper, req *messagev1.PushMessageReq) error {
	brief := &messagev1.BriefMessage{
		FromUser:    req.GetFromUser(),
		ToUser:      req.GetToUser(),
		ContentType: req.GetContentType(),
		Content:     req.GetContent(),
		MsgSeq:      req.GetMsgSeq(),
	}

	b, err := json.Marshal(brief)
	if err != nil {
		return err
	}

	return ww.WriteMessage(websocket.TextMessage, b)
}
