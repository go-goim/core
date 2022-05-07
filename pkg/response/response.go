package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"

	"github.com/yusank/goim/api/transport/response"
	"github.com/yusank/goim/pkg/log"
	"github.com/yusank/goim/pkg/mid"
)

func ErrorResp(c *gin.Context, err error) {
	log.Error("ErrorResp", "err", err)
	json(c, http.StatusOK, errorResp(err))
}

func errorResp(err error) *response.BaseResponse {
	switch t := err.(type) {
	case *response.BaseResponse:
		return t
	default:
		return response.ErrUnknown.SetMsg(err.Error())
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

func convertBodyToResponse(body interface{}) response.IResponse {
	var resp response.IResponse
	switch b := body.(type) {
	case response.IResponse:
		return b
	case proto.Message:
		resp = response.NewPbResponse(response.OK)
	default:
		resp = response.NewResponse(response.OK)
	}

	var err error
	resp, err = resp.SetData(body)
	if err != nil {
		return response.ErrUnknown.SetMsg(err.Error())
	}

	return resp
}

func json(c *gin.Context, code int, resp response.IResponse) {
	// get request id from context
	SetRequestID(c.GetString(mid.RequestIDKey))(resp.GetOrNewMeta())

	b, err := resp.Marshall()
	if err != nil {
		// this is a critical error
		log.Error("marshal response error.", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": response.ErrUnknown.GetCode(), "msg": err.Error()})
	}

	c.Data(code, jsonContentType, b)
}

// set meta data to response

type SetMetaFunc func(meta *response.Meta)

func SetPaging(page, size, total int) SetMetaFunc {
	return func(meta *response.Meta) {
		meta.SetPaging(page, size).SetTotal(total)
	}
}

func SetRequestID(requestID string) SetMetaFunc {
	return func(meta *response.Meta) {
		meta.SetRequestID(requestID)
	}
}

func SetExtra(key string, value string) SetMetaFunc {
	return func(meta *response.Meta) {
		meta.SetExtra(key, value)
	}
}

func SetExtraMap(m map[string]string) SetMetaFunc {
	return func(meta *response.Meta) {
		meta.SetExtraMap(m)
	}
}
