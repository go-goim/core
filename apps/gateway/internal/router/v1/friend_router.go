package v1

import (
	"github.com/gin-gonic/gin"

	responsepb "github.com/yusank/goim/api/transport/response"
	friendpb "github.com/yusank/goim/api/user/friend/v1"
	"github.com/yusank/goim/apps/gateway/internal/service"
	"github.com/yusank/goim/pkg/mid"
	"github.com/yusank/goim/pkg/request"
	"github.com/yusank/goim/pkg/response"
	"github.com/yusank/goim/pkg/router"
)

type FriendRouter struct {
	router.Router
}

func NewFriendRouter() *FriendRouter {
	return &FriendRouter{
		Router: &router.BaseRouter{},
	}
}

func (r *FriendRouter) Load(g *gin.RouterGroup) {
	g.GET("/list", r.listRelation)
	g.POST("/add-friend", r.addFriend)
	g.POST("/delete-friend", r.deleteFriend)
	g.POST("/accept-friend", r.acceptFriend)
	g.POST("/reject-friend", r.rejectFriend)
	g.POST("/block-friend", r.blockFriend)
	g.POST("/unblock-friend", r.unblockFriend)
}

func (r *FriendRouter) listRelation(c *gin.Context) {
	// no need to check uid
	uid := c.GetString("uid")
	req := &friendpb.QueryFriendListRequest{
		Uid: uid,
	}

	list, err := service.GetUserRelationService().ListUserRelation(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, list, response.SetTotal(len(list)))
}

func (r *FriendRouter) addFriend(c *gin.Context) {
	req := &friendpb.AddFriendRequest{}
	if err := c.ShouldBindWith(req, request.NonValidatePbJSONBinding); err != nil {
		response.ErrorResp(c, err)
		return
	}

	req.Uid = mid.GetUID(c)
	if err := req.Validate(); err != nil {
		response.ErrorResp(c, responsepb.NewBaseResponse(responsepb.Code_InvalidParams, err.Error()))
		return
	}

	result, err := service.GetUserRelationService().AddFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, result)
}

func (r *FriendRouter) deleteFriend(c *gin.Context) {
	req := &friendpb.BaseFriendRequest{}
	if err := c.ShouldBindWith(req, request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	err := service.GetUserRelationService().DeleteFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.OK(c)
}

func (r *FriendRouter) acceptFriend(c *gin.Context) {
	req := &friendpb.ConfirmFriendRequestReq{}
	if err := c.ShouldBindWith(req, request.NonValidatePbJSONBinding); err != nil {
		response.ErrorResp(c, err)
		return
	}

	err := service.GetUserRelationService().AcceptFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.OK(c)
}

func (r *FriendRouter) rejectFriend(c *gin.Context) {
	req := &friendpb.ConfirmFriendRequestReq{}
	if err := c.ShouldBindWith(req, request.NonValidatePbJSONBinding); err != nil {
		response.ErrorResp(c, err)
		return
	}

	err := service.GetUserRelationService().RejectFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.OK(c)
}

func (r *FriendRouter) blockFriend(c *gin.Context) {
	req := &friendpb.BaseFriendRequest{}
	if err := c.ShouldBindWith(req, request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	err := service.GetUserRelationService().BlockFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.OK(c)
}

func (r *FriendRouter) unblockFriend(c *gin.Context) {
	req := &friendpb.BaseFriendRequest{}
	if err := c.ShouldBindWith(req, request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	err := service.GetUserRelationService().UnblockFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.OK(c)
}
