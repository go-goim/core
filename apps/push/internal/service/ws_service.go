package service

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yusank/goim/pkg/log"

	"github.com/yusank/goim/apps/push/internal/app"
	"github.com/yusank/goim/apps/push/internal/data"
	"github.com/yusank/goim/pkg/conn/pool"
	"github.com/yusank/goim/pkg/conn/wrapper"
)

func HandleWsConn(ctx context.Context, c *websocket.Conn, uid string) {
	ww := wrapper.WrapWs(ctx, c, uid)
	ww.AddPingAction(func() error {
		return app.GetApplication().Redis.SetEX(context.Background(), data.GetUserOnlineAgentKey(uid), app.GetAgentID(), data.UserOnlineAgentKeyExpire).Err()
	})
	ww.AddCloseAction(func() error {
		ctx2, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		return app.GetApplication().Redis.Del(ctx2, data.GetUserOnlineAgentKey(uid)).Err()

	})

	go ww.Daemon()
	pool.Add(ww)

	ctx2, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	err := app.GetApplication().Redis.Set(ctx2, data.GetUserOnlineAgentKey(uid), app.GetAgentID(), data.UserOnlineAgentKeyExpire).Err()
	if err != nil {
		log.Error("redis set error", "key", data.GetUserOnlineAgentKey(uid), "error", err)
	}

	go func() {
		_ = c.WriteMessage(websocket.TextMessage, []byte("connect success"))
	}()
}
