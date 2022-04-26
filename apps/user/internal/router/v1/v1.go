package v1

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.RouterGroup) {
	NewUserRouter().Register(g.Group("/user"))
}
