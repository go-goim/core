package v1

import (
	"github.com/gin-gonic/gin"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/gateway/internal/service"
	"github.com/yusank/goim/pkg/mid"
	"github.com/yusank/goim/pkg/request"
	"github.com/yusank/goim/pkg/response"
	"github.com/yusank/goim/pkg/router"
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

func (r *MsgRouter) handleSendSingleUserMsg(c *gin.Context) {
	req := new(messagev1.SendMessageReq)
	if err := c.ShouldBindWith(req, request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	rsp, err := service.GetSendMessageService().SendMessage(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, rsp)
}

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

	response.SuccessResp(c, rsp)
}
