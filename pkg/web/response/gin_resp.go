package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-goim/api/errors"
	"github.com/go-goim/core/pkg/log"
	"github.com/go-goim/core/pkg/web"
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
	resp := NewResponseFromCode(errors.ErrorCode_OK).SetData(body)
	for _, f := range sf {
		f(resp.BaseResponse)
	}

	c.JSON(http.StatusOK, resp)
}

// OK is a shortcut for c.JSON(http.StatusOK, NewResponseFromPb(responsepb.OK))
func OK(c *gin.Context) {
	c.JSON(http.StatusOK, NewResponseFromCode(errors.ErrorCode_OK))
}

func errorResp(err error) *Response {
	return NewResponseFromPb(errors.NewErrorWithError(err))
}

type SetFunc func(resp *BaseResponse)

func SetPaging(paging *web.Paging) SetFunc {
	return func(resp *BaseResponse) {
		_ = resp.SetPaging(paging.Page, paging.PageSize)
	}
}

func SetTotal(total int32) SetFunc {
	return func(resp *BaseResponse) {
		_ = resp.SetTotal(total)
	}
}
