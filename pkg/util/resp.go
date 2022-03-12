package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
)

func ErrorResp(c *gin.Context, err error) {
	e := new(errors.Error)
	e.Code = errors.UnknownCode
	e.Message = err.Error()

	c.JSON(http.StatusOK, e)
}

func Success(c *gin.Context, body interface{}) {
	c.JSON(http.StatusOK, body)
}
