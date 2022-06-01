package v1

import (
	"github.com/gin-gonic/gin"

	responsepb "github.com/go-goim/core/api/transport/response"
	friendpb "github.com/go-goim/core/api/user/friend/v1"
	"github.com/go-goim/core/apps/gateway/internal/service"
	"github.com/go-goim/core/pkg/mid"
	"github.com/go-goim/core/pkg/request"
	"github.com/go-goim/core/pkg/response"
	"github.com/go-goim/core/pkg/router"
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

// @Summary 获取好友列表
// @Description 获取好友列表
// @Tags [gateway]好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Success 200 {object} response.Response{data=[]friendpb.Friend} "Success"
// @Failure 400 {object} response.Response{} "err"
// @Router /gateway/v1/user/friend/list [get]
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

// query user info before add friend
// @Summary 添加好友
// @Description 添加好友
// @Tags [gateway]好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body friendpb.AddFriendRequest true "body"
// @Success 200 {object} friendpb.AddFriendResult
// @Failure 400 {object} responsepb.BaseResponse "err"
// @Router /gateway/v1/user/friend/add-friend [post]
func (r *FriendRouter) addFriend(c *gin.Context) {
	req := &friendpb.AddFriendRequest{}
	if err := c.ShouldBindWith(req, request.NonValidatePbJSONBinding); err != nil {
		response.ErrorResp(c, err)
		return
	}

	req.Uid = mid.GetUID(c)
	if err := req.Validate(); err != nil {
		response.ErrorResp(c, responsepb.NewBaseResponseWithMessage(responsepb.Code_InvalidParams, err.Error()))
		return
	}

	result, err := service.GetUserRelationService().AddFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, result)
}

// @Summary 删除好友
// @Description 删除好友
// @Tags [gateway]好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body friendpb.BaseFriendRequest true "body"
// @Success 200 {object} responsepb.BaseResponse "Success"
// @Failure 400 {object} responsepb.BaseResponse "err"
// @Router /gateway/v1/user/friend/delete-friend [post]
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

// @Summary 接受好友请求
// @Description 接受好友请求
// @Tags [gateway]好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body friendpb.ConfirmFriendRequestReq true "body"
// @Success 200 {object} responsepb.BaseResponse "Success"
// @Failure 400 {object} responsepb.BaseResponse "err"
// @Router /gateway/v1/user/friend/accept-friend [post]
func (r *FriendRouter) acceptFriend(c *gin.Context) {
	req := &friendpb.ConfirmFriendRequestReq{}
	if err := c.ShouldBindWith(req, request.NonValidatePbJSONBinding); err != nil {
		response.ErrorResp(c, err)
		return
	}

	req.Uid = mid.GetUID(c)
	req.Action = friendpb.ConfirmFriendRequestAction_ACCEPT
	err := service.GetUserRelationService().AcceptFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.OK(c)
}

// @Summary 拒绝好友请求
// @Description 拒绝好友请求
// @Tags [gateway]好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body friendpb.ConfirmFriendRequestReq true "body"
// @Success 200 {object} responsepb.BaseResponse "Success"
// @Failure 400 {object} responsepb.BaseResponse "err"
// @Router /gateway/v1/user/friend/reject-friend [post]
func (r *FriendRouter) rejectFriend(c *gin.Context) {
	req := &friendpb.ConfirmFriendRequestReq{}
	if err := c.ShouldBindWith(req, request.NonValidatePbJSONBinding); err != nil {
		response.ErrorResp(c, err)
		return
	}

	req.Action = friendpb.ConfirmFriendRequestAction_REJECT
	req.Uid = mid.GetUID(c)
	err := service.GetUserRelationService().RejectFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.OK(c)
}

// @Summary 屏蔽好友
// @Description 屏蔽好友
// @Tags [gateway]好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body friendpb.BaseFriendRequest true "body"
// @Success 200 {object} responsepb.BaseResponse "Success"
// @Failure 400 {object} responsepb.BaseResponse "err"
// @Router /gateway/v1/user/friend/block-friend [post]
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

// @Summary 取消屏蔽好友
// @Description 取消屏蔽好友
// @Tags [gateway]好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body friendpb.BaseFriendRequest true "body"
// @Success 200 {object} responsepb.BaseResponse "Success"
// @Failure 400 {object} responsepb.BaseResponse "err"
// @Router /gateway/v1/user/friend/unblock-friend [post]
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
