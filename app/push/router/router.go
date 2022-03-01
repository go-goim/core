package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/yusank/goim/app/push/router/v1"
)

func RegisterRouter(g *gin.RouterGroup) {
	v1.Register(g.Group("/v1"))
}