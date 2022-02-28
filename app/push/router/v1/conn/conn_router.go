package conn

import (
	"github.com/gin-gonic/gin"
	"github.com/yusank/goim/pkg/conn"
)

func Register(g *gin.RouterGroup) {
	g.GET("/ws", WsConnHandler)
}

func WsConnHandler(c *gin.Context) {
	conn.WsHandler(c.Writer, c.Request)
}
