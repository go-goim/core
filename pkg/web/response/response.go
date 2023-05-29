package response

import (
	"github.com/go-goim/api/errors"
)

type BaseResponse struct {
	Code    int    `json:"code"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Meta    *Meta  `json:"meta,omitempty"`
}

func (r *BaseResponse) SetTotal(t int32) *BaseResponse {
	if r.Meta == nil {
		r.Meta = &Meta{}
	}

	r.Meta.Total = t
	return r
}

func (r *BaseResponse) SetPaging(page, size int32) *BaseResponse {
	if r.Meta == nil {
		r.Meta = &Meta{}
	}

	r.Meta.Page = page
	r.Meta.PageSize = size
	return r
}

func (r *BaseResponse) SetMsg(msg string) *BaseResponse {
	r.Message = msg
	return r
}

type Response struct {
	*BaseResponse `json:",inline"`
	Data          interface{} `json:"data,omitempty"`
}

func NewResponseFromPb(err *errors.Error) *Response {
	return &Response{
		BaseResponse: &BaseResponse{
			Code:    int(err.ErrorCode),
			Reason:  err.Reason,
			Message: err.Message,
		},
	}
}

func NewResponseFromCode(code errors.ErrorCode) *Response {
	return NewResponseFromPb(errors.NewErrorWithCode(code))
}

func (r *Response) SetData(data interface{}) *Response {
	r.Data = data

	return r
}

func (r *Response) SetMsg(msg string) *Response {
	r.Message = msg
	return r
}
