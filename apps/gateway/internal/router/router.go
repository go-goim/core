package router

import (
	"github.com/gin-gonic/gin"

	v1 "github.com/go-goim/goim/apps/gateway/internal/router/v1"
	"github.com/go-goim/goim/pkg/router"
)

type rootRouter struct {
	router.Router
}

func newRootRouter() *rootRouter {
	r := &rootRouter{
		Router: &router.BaseRouter{},
	}

	r.init()
	return r
}

func (r *rootRouter) init() {
	r.Register("/v1", v1.NewRouter())
}

func RegisterRouter(g *gin.RouterGroup) {
	r := newRootRouter()
	r.Load(g)
}
