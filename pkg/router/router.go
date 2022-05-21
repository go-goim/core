package router

import (
	"github.com/gin-gonic/gin"
)

// Router is interface for router
type Router interface {
	// Register registers routes.
	Register(path string, router Router)
	// Load registers routes to gin router in once.
	Load(group *gin.RouterGroup)
}
