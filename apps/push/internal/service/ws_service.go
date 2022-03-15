package service

import (
	"context"
	"net"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/websocket"

	"github.com/yusank/goim/apps/push/internal/app"
	"github.com/yusank/goim/apps/push/internal/data"
	"github.com/yusank/goim/pkg/conn"
)

func HandleWsConn(c *websocket.Conn, uid string) {
	wc := &conn.WsConn{
		Conn: c,
		Key:  uid,
	}

	c.SetPingHandler(func(message string) error {
		err := c.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(time.Second))
		if err == nil {
			log.Infof("get user:%s ping", wc.Key)
			app.GetApplication().Redis.SetEX(context.Background(), data.GetUserOnlineAgentKey(wc.Key), app.GetAgentID(), data.UserOnlineAgentKeyExpire).Err()
			return nil
		}

		if err == websocket.ErrCloseSent {
			return nil
		} else if e, ok := err.(net.Error); ok && e.Temporary() {
			return nil
		}
		return err
	})

	err := conn.PutConn(wc)
	if err != nil {
		log.Info(err)
	}

	err = app.GetApplication().Redis.Set(context.Background(), data.GetUserOnlineAgentKey(uid), app.GetAgentID(), data.UserOnlineAgentKeyExpire).Err()
	if err != nil {
		log.Info(err)
	}

	go func() {
		// write msg
		_ = c.SetWriteDeadline(time.Now().Add(time.Second))
		w, err1 := c.NextWriter(websocket.TextMessage)
		if err1 != nil {
			return
		}
		defer w.Close()

		_, err = w.Write([]byte("connect success"))
		if err != nil {
			return
		}
	}()
}
