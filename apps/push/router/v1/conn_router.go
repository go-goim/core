package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/yusank/goim/pkg/conn"
)

func wsConnHandler(c *gin.Context) {
	conn.WsHandler(c.Writer, c.Request)
}
