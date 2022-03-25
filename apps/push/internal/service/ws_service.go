package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/websocket"

	"github.com/yusank/goim/apps/push/internal/app"
	"github.com/yusank/goim/apps/push/internal/data"
	"github.com/yusank/goim/pkg/pool"
	goimwebsocket "github.com/yusank/goim/pkg/websocket"
)

func HandleWsConn(c *websocket.Conn, uid string) {
	ww := goimwebsocket.WrapWs(c, uid)
	ww.AddPingAction(func() error {
		return app.GetApplication().Redis.SetEX(context.Background(), data.GetUserOnlineAgentKey(uid), app.GetAgentID(), data.UserOnlineAgentKeyExpire).Err()
	})
	ww.AddCloseAction(func() error {
		return app.GetApplication().Redis.Del(context.Background(), data.GetUserOnlineAgentKey(uid)).Err()

	})

	pool.Add(ww)

	err := app.GetApplication().Redis.Set(context.Background(), data.GetUserOnlineAgentKey(uid), app.GetAgentID(), data.UserOnlineAgentKeyExpire).Err()
	if err != nil {
		log.Info(err)
	}

	go func() {
		_ = c.WriteMessage(websocket.TextMessage, []byte("connect success"))
	}()
}
