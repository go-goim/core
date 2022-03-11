package service

import (
	"context"
	"encoding/json"
	"net"
	"time"

	"github.com/gorilla/websocket"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/push/internal/app"
	"github.com/yusank/goim/apps/push/internal/data"
	"github.com/yusank/goim/pkg/conn"
)

type WsConn struct {
	*websocket.Conn
	key      string
	pingChan chan struct{}
}

func HandleWsConn(c *websocket.Conn, uid string) {
	wc := &WsConn{
		Conn:     c,
		key:      uid,
		pingChan: make(chan struct{}, 1),
	}

	wc.SetPingHandler(wc.pingFunc)
	_ = conn.PutConn(wc)

	_ = app.GetApplication().Redis.Set(context.Background(), data.GetUserOnlineAgentKey(uid), app.GetAgentID(), data.UserOnlineAgentKeyExpire).Err()

	// write msg
	_ = c.SetWriteDeadline(time.Now().Add(time.Second))
	w, err := wc.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}

	_, err = w.Write([]byte("connect success"))
	if err != nil {
		return
	}

	_ = w.Close()
}

func (wc *WsConn) PushMessage(message *messagev1.PushMessageReq) error {
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_ = wc.SetWriteDeadline(time.Now().Add(time.Second))
	w, err := wc.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return w.Close()
}

func (w *WsConn) Key() string {
	return w.key
}

func (w *WsConn) Ping() <-chan struct{} {
	return w.pingChan
}

func (w *WsConn) pingFunc(message string) error {
	err := w.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(time.Second))
	if err == nil {
		select {
		// try put ping
		case w.pingChan <- struct{}{}:
		// non-blocking
		default:
		}

		return nil
	}

	if err == websocket.ErrCloseSent {
		return nil
	} else if e, ok := err.(net.Error); ok && e.Temporary() {
		return nil
	}
	return err
}
