package v1

import (
	"github.com/gin-gonic/gin"

	userv1 "github.com/yusank/goim/api/user/v1"
	"github.com/yusank/goim/apps/gateway/internal/service"
	"github.com/yusank/goim/pkg/mid"
	"github.com/yusank/goim/pkg/router"
	"github.com/yusank/goim/pkg/util"
)

type UserRelationRouter struct {
	router.Router
}

func NewUserRelationRouter() *UserRelationRouter {
	return &UserRelationRouter{
		Router: &router.BaseRouter{},
	}
}

func (r *UserRelationRouter) Load(g *gin.RouterGroup) {
	g.GET("/list", r.listRelation)
	g.POST("/add-friend", r.addFriend)
	g.POST("/remove-friend", r.removeFriend)
	g.POST("/accept-friend", r.acceptFriend)
	g.POST("/reject-friend", r.rejectFriend)
	g.POST("/block-friend", r.blockFriend)
	g.POST("/unblock-friend", r.unblockFriend)
}

func (r *UserRelationRouter) listRelation(c *gin.Context) {
	// no need to check uid
	uid := c.GetString("uid")

	resp, err := service.GetUserRelationService().ListUserRelation(mid.GetContext(c), &userv1.QueryUserRelationListRequest{
		Uid: uid,
	})
	if err != nil {
		util.ErrorResp(c, err)
		return
	}

	util.Success(c, resp)
}

func (r *UserRelationRouter) addFriend(c *gin.Context) {
	req := &userv1.AddFriendRequest{}
	if err := c.ShouldBindWith(req, util.PbJSONBinding{}); err != nil {
		util.ErrorResp(c, err)
		return
	}

	resp, err := service.GetUserRelationService().AddFriend(mid.GetContext(c), req)
	if err != nil {
		util.ErrorResp(c, err)
		return
	}

	util.Success(c, resp)
}

func (r *UserRelationRouter) removeFriend(c *gin.Context) {
	req := &userv1.RemoveFriendRequest{}
	if err := c.ShouldBindWith(req, util.PbJSONBinding{}); err != nil {
		util.ErrorResp(c, err)
		return
	}

	err := service.GetUserRelationService().RemoveFriend(mid.GetContext(c), req)
	if err != nil {
		util.ErrorResp(c, err)
		return
	}

	util.Success(c, gin.H{"code": 0})
}

func (r *UserRelationRouter) acceptFriend(c *gin.Context) {
	req := &userv1.AcceptFriendRequest{}
	if err := c.ShouldBindWith(req, util.PbJSONBinding{}); err != nil {
		util.ErrorResp(c, err)
		return
	}

	err := service.GetUserRelationService().AcceptFriend(mid.GetContext(c), req)
	if err != nil {
		util.ErrorResp(c, err)
		return
	}

	util.Success(c, gin.H{"code": 0})
}

func (r *UserRelationRouter) rejectFriend(c *gin.Context) {
	req := &userv1.RejectFriendRequest{}
	if err := c.ShouldBindWith(req, util.PbJSONBinding{}); err != nil {
		util.ErrorResp(c, err)
		return
	}

	err := service.GetUserRelationService().RejectFriend(mid.GetContext(c), req)
	if err != nil {
		util.ErrorResp(c, err)
		return
	}

	util.Success(c, gin.H{"code": 0})
}

func (r *UserRelationRouter) blockFriend(c *gin.Context) {
	req := &userv1.BlockFriendRequest{}
	if err := c.ShouldBindWith(req, util.PbJSONBinding{}); err != nil {
		util.ErrorResp(c, err)
		return
	}

	err := service.GetUserRelationService().BlockFriend(mid.GetContext(c), req)
	if err != nil {
		util.ErrorResp(c, err)
		return
	}

	util.Success(c, gin.H{"code": 0})
}

func (r *UserRelationRouter) unblockFriend(c *gin.Context) {
	req := &userv1.UnblockFriendRequest{}
	if err := c.ShouldBindWith(req, util.PbJSONBinding{}); err != nil {
		util.ErrorResp(c, err)
		return
	}

	err := service.GetUserRelationService().UnblockFriend(mid.GetContext(c), req)
	if err != nil {
		util.ErrorResp(c, err)
		return
	}

	util.Success(c, gin.H{"code": 0})
}
