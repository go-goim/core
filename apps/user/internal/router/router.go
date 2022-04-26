package router

import (
	"github.com/gin-gonic/gin"

	v1 "github.com/yusank/goim/apps/user/internal/router/v1"
)

func RegisterRouters(g *gin.RouterGroup) {
	v1.RegisterRoutes(g.Group("/v1"))
}
