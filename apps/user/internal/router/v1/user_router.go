package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"

	userv1 "github.com/yusank/goim/api/user/v1"
	"github.com/yusank/goim/apps/user/internal/service"
	"github.com/yusank/goim/pkg/mid"
	"github.com/yusank/goim/pkg/util"
)

type UserRouter struct {
}

func NewUserRouter() *UserRouter {
	return &UserRouter{}
}

func (r *UserRouter) Register(router *gin.RouterGroup) {
	router.GET("", r.GetUser)
}

func (r *UserRouter) GetUser(c *gin.Context) {
	// get uid from query
	uid := c.Query("uid")
	if uid == "" {
		util.ErrorResp(c, fmt.Errorf("uid is empty"))
		return
	}

	// get user info from service
	req := &userv1.GetUserRequest{Uid: uid}
	resp, err := service.GetUserService().GetUser(mid.GetContext(c), req)
	if err != nil {
		util.ErrorResp(c, err)
		return
	}

	util.Success(c, resp)
}
