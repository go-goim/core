package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func ErrorResp(c *gin.Context, err error) {
	e := new(errors.Error)
	e.Code = errors.UnknownCode
	e.Message = err.Error()

	c.JSON(http.StatusOK, e)
}

const jsonContentType = "application/json; charset=utf-8"

func Success(c *gin.Context, body interface{}) {
	v, ok := body.(proto.Message)
	if ok {
		o := protojson.MarshalOptions{EmitUnpopulated: true}
		b, _ := o.Marshal(v)
		c.Data(http.StatusOK, jsonContentType, b)
		return
	}

	c.JSON(http.StatusOK, body)
}
