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
			_ = app.GetApplication().Redis.SetEX(context.Background(), data.GetUserOnlineAgentKey(wc.Key), app.GetAgentID(), data.UserOnlineAgentKeyExpire).Err()
			return nil
		}

		if err == websocket.ErrCloseSent {
			return nil
		} else if e, ok := err.(net.Error); ok && e.Temporary() {
			return nil
		}
		return err
	})

	c.SetCloseHandler(func(code int, text string) error {
		message := websocket.FormatCloseMessage(code, "")
		_ = c.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		_ = app.GetApplication().Redis.Del(context.Background(), data.GetUserOnlineAgentKey(wc.Key)).Err()
		return nil
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
		_ = c.WriteMessage(websocket.TextMessage, []byte("connect success"))
	}()
}
