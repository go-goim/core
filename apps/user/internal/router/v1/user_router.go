package v1

import (
	"github.com/gin-gonic/gin"

	userv1 "github.com/yusank/goim/api/user/v1"
	"github.com/yusank/goim/apps/user/internal/service"
	"github.com/yusank/goim/pkg/mid"
	"github.com/yusank/goim/pkg/resp"
	"github.com/yusank/goim/pkg/router"
)

type UserRouter struct {
	router.Router
}

func NewUserRouter() *UserRouter {
	return &UserRouter{
		Router: &router.BaseRouter{},
	}
}

func (r *UserRouter) Load(router *gin.RouterGroup) {
	router.GET("", r.GetUser)
}

func (r *UserRouter) GetUser(c *gin.Context) {
	// get uid from query
	req := &userv1.GetUserInfoRequest{
		Uid: c.Query("uid"),
	}
	if err := req.ValidateAll(); err != nil {
		resp.ErrorResp(c, err)
		return
	}

	// get user info from service

	result, err := service.GetUserService().GetUser(mid.GetContext(c), req)
	if err != nil {
		resp.ErrorResp(c, err)
		return
	}

	resp.SuccessResp(c, result)
}
