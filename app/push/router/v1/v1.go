package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/yusank/goim/app/push/router/v1/conn"
)

func Register(g *gin.RouterGroup) {
	// todo add register
	conn.Register(g.Group("/conn"))
}
