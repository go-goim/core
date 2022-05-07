package v1

import (
	"github.com/gin-gonic/gin"

	relationv1 "github.com/yusank/goim/api/user/relation/v1"
	"github.com/yusank/goim/apps/gateway/internal/service"
	"github.com/yusank/goim/pkg/mid"
	"github.com/yusank/goim/pkg/request"
	"github.com/yusank/goim/pkg/response"
	"github.com/yusank/goim/pkg/router"
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
	g.POST("/delete-friend", r.deleteFriend)
	g.POST("/accept-friend", r.acceptFriend)
	g.POST("/reject-friend", r.rejectFriend)
	g.POST("/block-friend", r.blockFriend)
	g.POST("/unblock-friend", r.unblockFriend)
}

func (r *UserRelationRouter) listRelation(c *gin.Context) {
	// no need to check uid
	uid := c.GetString("uid")
	req := &relationv1.QueryUserRelationListRequest{
		Uid: uid,
	}

	list, err := service.GetUserRelationService().ListUserRelation(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	total := len(list.UserRelationList)
	response.SuccessResp(c, list, response.SetPaging(1, total, total))
}

func (r *UserRelationRouter) addFriend(c *gin.Context) {
	req := &relationv1.AddFriendRequest{}
	if err := c.ShouldBindWith(req, request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	result, err := service.GetUserRelationService().AddFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, result)
}

func (r *UserRelationRouter) deleteFriend(c *gin.Context) {
	req := &relationv1.RelationPair{}
	if err := c.ShouldBindWith(req, request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	rsp, err := service.GetUserRelationService().DeleteFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, rsp)
}

func (r *UserRelationRouter) acceptFriend(c *gin.Context) {
	req := &relationv1.RelationPair{}
	if err := c.ShouldBindWith(req, request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	rsp, err := service.GetUserRelationService().AcceptFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, rsp)
}

func (r *UserRelationRouter) rejectFriend(c *gin.Context) {
	req := &relationv1.RelationPair{}
	if err := c.ShouldBindWith(req, request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	rsp, err := service.GetUserRelationService().RejectFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, rsp)
}

func (r *UserRelationRouter) blockFriend(c *gin.Context) {
	req := &relationv1.RelationPair{}
	if err := c.ShouldBindWith(req, request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	rsp, err := service.GetUserRelationService().BlockFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, rsp)
}

func (r *UserRelationRouter) unblockFriend(c *gin.Context) {
	req := &relationv1.RelationPair{}
	if err := c.ShouldBindWith(req, request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	rsp, err := service.GetUserRelationService().UnblockFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, rsp)
}
