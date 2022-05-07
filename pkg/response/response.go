package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"

	transportv1 "github.com/yusank/goim/api/transport/v1"
	"github.com/yusank/goim/pkg/log"
	"github.com/yusank/goim/pkg/mid"
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

func SuccessResp(c *gin.Context, body interface{}, setFunc ...SetMetaFunc) {
	resp := convertBodyToResponse(body)

	meta := resp.GetOrNewMeta()
	for _, f := range setFunc {
		f(meta)
	}
	resp.SetMeta(meta)

	json(c, http.StatusOK, resp)
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

func json(c *gin.Context, code int, resp transportv1.IResponse) {
	// get request id from context
	SetRequestID(c.GetString(mid.RequestIDKey))(resp.GetOrNewMeta())

	b, err := resp.Marshall()
	if err != nil {
		// this is a critical error
		log.Error("marshal response error.", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": transportv1.ResponseUnknownError.GetCode(), "msg": err.Error()})
	}

	c.Data(code, jsonContentType, b)
}

// set meta data to response

type SetMetaFunc func(meta *transportv1.Meta)

func SetPaging(page, size, total int) SetMetaFunc {
	return func(meta *transportv1.Meta) {
		meta.SetPaging(page, size).SetTotal(total)
	}
}

func SetRequestID(requestID string) SetMetaFunc {
	return func(meta *transportv1.Meta) {
		meta.SetRequestID(requestID)
	}
}

func SetExtra(key string, value string) SetMetaFunc {
	return func(meta *transportv1.Meta) {
		meta.SetExtra(key, value)
	}
}

func SetExtraMap(m map[string]string) SetMetaFunc {
	return func(meta *transportv1.Meta) {
		meta.SetExtraMap(m)
	}
}
