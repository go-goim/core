package service

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"

	responsepb "github.com/go-goim/api/transport/response"

	"github.com/go-goim/core/pkg/log"

	messagev1 "github.com/go-goim/api/message/v1"

	"github.com/go-goim/core/pkg/conn/pool"
	"github.com/go-goim/core/pkg/conn/wrapper"
	"github.com/go-goim/core/pkg/graceful"
	"github.com/go-goim/core/pkg/worker"
)

type PushMessager struct {
	messagev1.UnimplementedPushMessagerServer
	workerPool *worker.Pool
}

var (
	pm     *PushMessager
	pmOnce sync.Once
)

func GetPushMessager() *PushMessager {
	pmOnce.Do(func() {
		pm = new(PushMessager)
		pm.workerPool = worker.NewPool(100, 20)
		graceful.Register(pm.workerPool.Shutdown)
	})

	return pm
}

func (p *PushMessager) PushMessage(ctx context.Context, req *messagev1.PushMessageReq) (resp *responsepb.BaseResponse, err error) {
	log.Info("PUSH receive msg|", req.String())
	resp = responsepb.Code_OK.BaseResponse()
	if req.GetPushMessageType() == messagev1.PushMessageType_Broadcast {
		// cannot use request ctx in async function.It may kill the goroutine after this request finished.
		go p.handleBroadcastAsync(context.Background(), req)
		return
	}

	c := pool.Get(req.GetToUser())
	if c == nil {
		log.Info("PUSH| user conn not found=", req.GetToUser())
		resp = responsepb.Code_UserNotOnline.BaseResponse()
		return
	}

	err1 := PushMessage(c.(*wrapper.WebsocketWrapper), req)
	if err1 == nil {
		return
	}

	log.Info("PUSH| push err=", err1)
	resp = responsepb.NewBaseResponseWithError(err1)
	return
}

func (p *PushMessager) handleBroadcastAsync(ctx context.Context, req *messagev1.PushMessageReq) {
	ch := pool.LoadAllConn()
	wf := func() error {
		for c := range ch {
			select {
			case <-c.Done():
				continue
			default:
				if c.Err() != nil {
					continue
				}
			}

			ww := c.(*wrapper.WebsocketWrapper)
			if err := PushMessage(ww, req); err != nil {
				log.Info("PushMessage err=", err)
			}
		}

		return nil
	}

	result := p.workerPool.Submit(ctx, wf, 5)
	log.Info("PUSH| broadcast result=", result, "| status=", result.Status(), "| err=", result.Err())
	if result.Status() == worker.TaskStatusQueueFull {
		log.Info("worker queue buffer is full, should set more buffer")
	}
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
