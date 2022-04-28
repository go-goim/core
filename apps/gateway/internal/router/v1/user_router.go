package v1

import (
	"github.com/gin-gonic/gin"

	userv1 "github.com/yusank/goim/api/user/v1"
	"github.com/yusank/goim/apps/gateway/internal/service"
	"github.com/yusank/goim/pkg/mid"
	"github.com/yusank/goim/pkg/util"
)

type UserRouter struct {
}

func NewUserRouter() *UserRouter {
	return &UserRouter{}
}

func (r *UserRouter) Register(router *gin.RouterGroup) {
	router.POST("/login", r.login)
	router.POST("/register", r.register)
	router.POST("/update", r.updateUserInfo)
}

func (r *UserRouter) login(c *gin.Context) {
	var req = &userv1.UserLoginRequest{}
	if err := c.ShouldBind(req); err != nil {
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

	util.Success(c, gin.H{"user": user})
}

func (r *UserRouter) register(c *gin.Context) {
	var req = &userv1.CreateUserRequest{}
	if err := c.ShouldBind(req); err != nil {
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

	util.Success(c, gin.H{"user": user})
}

func (r *UserRouter) updateUserInfo(c *gin.Context) {
	var req = &userv1.UpdateUserRequest{}
	if err := c.ShouldBind(req); err != nil {
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

	util.Success(c, gin.H{"user": user})
}
