package v1

import (
	"github.com/gin-gonic/gin"

	userv1 "github.com/yusank/goim/api/user/v1"
	"github.com/yusank/goim/apps/gateway/internal/service"
	"github.com/yusank/goim/pkg/mid"
	"github.com/yusank/goim/pkg/request"
	"github.com/yusank/goim/pkg/response"
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
	friend := NewFriendRouter()
	friend.Load(router.Group("/friend", mid.AuthJwtCookie))

	router.POST("/login", r.login)
	router.POST("/register", r.register)
	router.POST("/update", mid.AuthJwtCookie, r.updateUserInfo)
}

func (r *UserRouter) login(c *gin.Context) {
	var req = &userv1.UserLoginRequest{}
	if err := c.ShouldBindWith(req, &request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	user, err := service.GetUserService().Login(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	if err = mid.SetJwtToHeader(c, user.Uid); err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, user)
}

func (r *UserRouter) register(c *gin.Context) {
	var req = &userv1.CreateUserRequest{}
	if err := c.ShouldBindWith(req, &request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	user, err := service.GetUserService().Register(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, user)
}

func (r *UserRouter) updateUserInfo(c *gin.Context) {
	var req = &userv1.UpdateUserRequest{}
	if err := c.ShouldBindWith(req, &request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	user, err := service.GetUserService().UpdateUser(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, user)
}
