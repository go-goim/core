package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/yusank/goim/apps/push/internal/service"
	"github.com/yusank/goim/pkg/mid"
)

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1 << 16,
	ReadBufferSize:  1024,
}

func wsConnHandler(c *gin.Context) {
	// todo use check uid/token middleware before this handler
	uid := c.GetHeader("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "uid not found"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	service.HandleWsConn(mid.GetContext(c), conn, uid)
}
