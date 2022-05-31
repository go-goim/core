package v1

import (
	"github.com/gin-gonic/gin"

	messagev1 "github.com/go-goim/goim/api/message/v1"
	responsepb "github.com/go-goim/goim/api/transport/response"
	"github.com/go-goim/goim/apps/gateway/internal/service"
	"github.com/go-goim/goim/pkg/mid"
	"github.com/go-goim/goim/pkg/request"
	"github.com/go-goim/goim/pkg/response"
	"github.com/go-goim/goim/pkg/router"
)

type MsgRouter struct {
	router.Router
}

func NewMsgRouter() *MsgRouter {
	return &MsgRouter{
		Router: &router.BaseRouter{},
	}
}

func (r *MsgRouter) Load(g *gin.RouterGroup) {
	g.Use(mid.AuthJwtCookie)
	offline := NewOfflineMessageRouter()
	offline.Load(g.Group("/offline_msg"))

	g.POST("/send_msg", r.handleSendSingleUserMsg)
	g.POST("/broadcast", r.handleSendBroadcastMsg)
}

// @Summary 发送单聊消息
// @Description 发送单聊消息
// @Tags [gateway]message
// @Accept  json
// @Produce  json
// @Param   Authorization header string true "token"
// @Param   req body messagev1.SendMessageReq true "req"
// @Success 200 {object} messagev1.SendMessageResp
// @Failure 200 {object} response.Response
// @Failure 401 {null} null
// @Router /gateway/v1/message/send_msg [post]
func (r *MsgRouter) handleSendSingleUserMsg(c *gin.Context) {
	req := new(messagev1.SendMessageReq)
	if err := c.ShouldBindWith(req, request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ErrorResp(c, responsepb.NewBaseResponseWithMessage(responsepb.Code_InvalidParams, err.Error()))
		return
	}

	rsp, err := service.GetSendMessageService().SendMessage(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, gin.H{
		"msg_seq": rsp.GetMsgSeq(),
	})
}

// @Summary 发送广播消息
// @Description 发送广播消息
// @Tags [gateway]message
// @Accept  json
// @Produce  json
// @Param   Authorization header string true "token"
// @Param   req body messagev1.SendMessageReq true "req"
// @Success 200 {object} messagev1.SendMessageResp
// @Failure 200 {object} response.Response
// @Failure 401 {null} null
// @Router /gateway/v1/message/broadcast [post]
func (r *MsgRouter) handleSendBroadcastMsg(c *gin.Context) {
	req := new(messagev1.SendMessageReq)
	if err := c.ShouldBindWith(req, request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	rsp, err := service.GetSendMessageService().Broadcast(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, gin.H{
		"msg_seq": rsp.GetMsgSeq(),
	})
}
