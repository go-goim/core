package router

import (
	"github.com/gin-gonic/gin"
)

// BaseRouter is a base router to be embedded in all routers, and it implements Router interface.
type BaseRouter struct {
	path        string
	childRoutes map[string]Router
}

func (b *BaseRouter) Register(path string, router Router) {
	b.path = path
	if b.childRoutes == nil {
		b.childRoutes = make(map[string]Router)
	}

	b.childRoutes[path] = router
}

func (b *BaseRouter) Load(rg *gin.RouterGroup) {
	for pattern, router := range b.childRoutes {
		router.Load(rg.Group(pattern))
	}
}
