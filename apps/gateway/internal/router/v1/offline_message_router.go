package v1

import (
	"github.com/gin-gonic/gin"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/gateway/internal/service"
	"github.com/yusank/goim/pkg/mid"
	"github.com/yusank/goim/pkg/util"
)

func handleQueryOfflineMessage(c *gin.Context) {
	req := new(messagev1.QueryOfflineMessageReq)
	if err := c.ShouldBindJSON(req); err != nil {
		util.ErrorResp(c, err)
		return
	}

	rsp, err := service.GetOfflineMessageService().QueryOfflineMsg(mid.GetContext(c), req)
	if err != nil {
		util.ErrorResp(c, err)
		return
	}

	util.Success(c, rsp)
}
