package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/protobuf/proto"

	"github.com/yusank/goim/pkg/log"
)

func ErrorResp(c *gin.Context, err error) {
	e := new(errors.Error)
	e.Code = errors.UnknownCode
	e.Message = err.Error()

	log.Debug("ErrorResp", "err", err, "route", c.Request.URL.Path)

	c.JSON(http.StatusOK, e)
}

func ErrorRespWithStatus(c *gin.Context, httpCode int, err error) {
	e := new(errors.Error)
	e.Code = errors.UnknownCode
	e.Message = err.Error()

	log.Debug("ErrorResp", "err", err, "route", c.Request.URL.Path)
	c.JSON(httpCode, e)
}

const jsonContentType = "application/json; charset=utf-8"

func SuccessResp(c *gin.Context, body interface{}) {
	v, ok := body.(proto.Message)
	if ok {
		b, _ := MarshallPb(v)
		c.Data(http.StatusOK, jsonContentType, b)
		return
	}

	c.JSON(http.StatusOK, body)
}
