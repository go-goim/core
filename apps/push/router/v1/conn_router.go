package v1

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/yusank/goim/apps/push/service"
)

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1 << 16,
	ReadBufferSize:  1024,
}

func wsConnHandler(c *gin.Context) {
	//todo use check uid/token middleware before this handler
	uid := c.GetHeader("uid")
	if uid == "" {
		log.Println("uid not found")
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// todo
		log.Println(err)
		return
	}

	service.HandleWsConn(conn, uid)
}
