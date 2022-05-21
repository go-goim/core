package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	responsepb "github.com/yusank/goim/api/transport/response"
	"github.com/yusank/goim/pkg/log"
)

func ErrorResp(c *gin.Context, err error) {
	log.Error("ErrorResp", "err", err)
	c.JSON(http.StatusOK, errorResp(err))
}

func ErrorRespWithStatus(c *gin.Context, httpCode int, err error) {
	log.Error("ErrorResp", "err", err)
	c.JSON(httpCode, errorResp(err))
}

func SuccessResp(c *gin.Context, body interface{}, sf ...SetFunc) {
	resp := NewResponseFromCode(responsepb.Code_OK).SetData(body)
	for _, f := range sf {
		f(resp.BaseResponse)
	}

	c.JSON(http.StatusOK, resp)
}

// OK is a shortcut for c.JSON(http.StatusOK, NewResponseFromPb(responsepb.OK))
func OK(c *gin.Context) {
	c.JSON(http.StatusOK, NewResponseFromCode(responsepb.Code_OK))
}

func errorResp(err error) *Response {
	switch t := err.(type) {
	case *responsepb.BaseResponse:
		return NewResponseFromPb(t)
	default:
		return NewResponseFromPb(responsepb.NewBaseResponseWithError(err))
	}
}

type SetFunc func(resp *BaseResponse)

func SetPaging(page, size int) SetFunc {
	return func(resp *BaseResponse) {
		_ = resp.SetPaging(page, size)
	}
}

func SetTotal(total int) SetFunc {
	return func(resp *BaseResponse) {
		_ = resp.SetTotal(total)
	}
}
