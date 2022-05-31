package v1

import (
	"github.com/go-goim/goim/pkg/router"
)

type Router struct {
	router.Router
}

func NewRouter() *Router {
	r := &Router{
		Router: &router.BaseRouter{},
	}
	r.Init()
	return r
}

func (r *Router) Init() {
	r.Register("/user", NewUserRouter())
	r.Register("/msg", NewMsgRouter())
	r.Register("/discovery", NewDiscoverRouter())
}
