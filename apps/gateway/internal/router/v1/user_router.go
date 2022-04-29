package v1

import (
	"github.com/gin-gonic/gin"

	userv1 "github.com/yusank/goim/api/user/v1"
	"github.com/yusank/goim/apps/gateway/internal/service"
	"github.com/yusank/goim/pkg/mid"
	"github.com/yusank/goim/pkg/router"
	"github.com/yusank/goim/pkg/util"
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
	router.POST("/login", r.login)
	router.POST("/register", r.register)
	router.POST("/update", r.updateUserInfo)
}

func (r *UserRouter) login(c *gin.Context) {
	var req = &userv1.UserLoginRequest{}
	if err := c.ShouldBindWith(req, &util.PbJsonBinding{}); err != nil {
		util.ErrorResp(c, err)
		return
	}

	if err := req.ValidateAll(); err != nil {
		util.ErrorResp(c, err)
		return
	}

	user, err := service.GetUserService().Login(mid.GetContext(c), req)
	if err != nil {
		util.ErrorResp(c, err)
		return
	}

	util.Success(c, user)
}

func (r *UserRouter) register(c *gin.Context) {
	var req = &userv1.CreateUserRequest{}
	if err := c.ShouldBindWith(req, &util.PbJsonBinding{}); err != nil {
		util.ErrorResp(c, err)
		return
	}

	if err := req.ValidateAll(); err != nil {
		util.ErrorResp(c, err)
		return
	}

	user, err := service.GetUserService().Register(mid.GetContext(c), req)
	if err != nil {
		util.ErrorResp(c, err)
		return
	}

	util.Success(c, user)
}

func (r *UserRouter) updateUserInfo(c *gin.Context) {
	var req = &userv1.UpdateUserRequest{}
	if err := c.ShouldBindWith(req, &util.PbJsonBinding{}); err != nil {
		util.ErrorResp(c, err)
		return
	}

	if err := req.ValidateAll(); err != nil {
		util.ErrorResp(c, err)
		return
	}

	user, err := service.GetUserService().UpdateUser(mid.GetContext(c), req)
	if err != nil {
		util.ErrorResp(c, err)
		return
	}

	util.Success(c, user)
}
