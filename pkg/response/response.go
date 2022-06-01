package response

import (
	responsepb "github.com/go-goim/core/api/transport/response"
)

type BaseResponse struct {
	Code    int    `json:"code"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Meta    *Meta  `json:"meta,omitempty"`
}

func (r *BaseResponse) SetTotal(t int) *BaseResponse {
	if r.Meta == nil {
		r.Meta = &Meta{}
	}

	r.Meta.Total = t
	return r
}

func (r *BaseResponse) SetPaging(page, size int) *BaseResponse {
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
	*BaseResponse
	Data interface{} `json:"data,omitempty"`
}

func NewResponseFromPb(base *responsepb.BaseResponse) *Response {
	return &Response{
		BaseResponse: &BaseResponse{
			Code:    int(base.Code),
			Reason:  base.Reason,
			Message: base.Message,
		},
	}
}

func NewResponseFromCode(code responsepb.Code) *Response {
	return NewResponseFromPb(code.BaseResponse())
}

func (r *Response) SetData(data interface{}) *Response {
	r.Data = data

	return r
}

func (r *Response) SetMsg(msg string) *Response {
	r.Message = msg
	return r
}
