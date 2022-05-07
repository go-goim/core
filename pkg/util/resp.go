package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"

	transportv1 "github.com/yusank/goim/api/transport/v1"
	"github.com/yusank/goim/pkg/log"
)

func ErrorResp(c *gin.Context, err error) {
	log.Error("ErrorResp", "err", err)
	json(c, http.StatusOK, errorResp(err))
}

func errorResp(err error) *transportv1.BaseResponse {
	switch t := err.(type) {
	case *transportv1.BaseResponse:
		return t
	default:
		return transportv1.ResponseUnknownError.SetMsg(err.Error())
	}
}

func ErrorRespWithStatus(c *gin.Context, httpCode int, err error) {
	log.Error("ErrorResp", "err", err)
	json(c, httpCode, errorResp(err))
}

const jsonContentType = "application/json; charset=utf-8"

func SuccessResp(c *gin.Context, body interface{}) {
	json(c, http.StatusOK, body)
}

func convertBodyToResponse(body interface{}) transportv1.IResponse {
	var resp transportv1.IResponse
	switch b := body.(type) {
	case transportv1.IResponse:
		return b
	case proto.Message:
		resp = transportv1.NewPbResponse(transportv1.ResponseOK)
	default:
		resp = transportv1.NewResponse(transportv1.ResponseOK)
	}

	var err error
	resp, err = resp.SetData(body)
	if err != nil {
		return transportv1.ResponseUnknownError.SetMsg(err.Error())
	}

	return resp
}

func json(c *gin.Context, code int, body interface{}) {
	resp := convertBodyToResponse(body)
	// TODO set meta data
	b, err := resp.Marshall()
	if err != nil {
		// this is a critical error
		log.Error("marshal response error.", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": transportv1.ResponseUnknownError.GetCode(), "msg": err.Error()})
	}

	c.Data(code, jsonContentType, b)
}
