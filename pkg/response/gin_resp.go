package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	responsepb "github.com/yusank/goim/api/transport/response"
	"github.com/yusank/goim/pkg/log"
)

func ErrorResp(c *gin.Context, err error) {
	log.Error("ErrorResp", "err", err)
	c.JSON(http.StatusOK, NewResponseFromPb(errorResp(err)))
}

func ErrorRespWithStatus(c *gin.Context, httpCode int, err error) {
	log.Error("ErrorResp", "err", err)
	c.JSON(httpCode, NewResponseFromPb(errorResp(err)))
}

func SuccessResp(c *gin.Context, baseResp *responsepb.BaseResponse, body interface{}) {
	resp := NewResponseFromPb(baseResp)
	if body != nil {
		resp.SetData(body)
	}

	c.JSON(http.StatusOK, resp)
}

// OK is a shortcut for c.JSON(http.StatusOK, NewResponseFromPb(responsepb.OK))
func OK(c *gin.Context) {
	c.JSON(http.StatusOK, NewResponseFromPb(responsepb.OK))
}

func SuccessWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, NewResponseFromPb(responsepb.OK).SetData(data))
}

func SuccessWithBaseResp(c *gin.Context, baseResp *responsepb.BaseResponse) {
	c.JSON(http.StatusOK, NewResponseFromPb(baseResp))
}

func errorResp(err error) *responsepb.BaseResponse {
	switch t := err.(type) {
	case *responsepb.BaseResponse:
		return t
	default:
		return responsepb.ErrUnknown.SetMsg(err.Error())
	}
}
