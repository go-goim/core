package response

import (
	responsepb "github.com/yusank/goim/api/transport/response"
)

type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Meta *Meta  `json:"meta,omitempty"`
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
	r.Msg = msg
	return r
}

type Response struct {
	*BaseResponse
	Data interface{} `json:"data,omitempty"`
}

func NewResponse(code int, msg string) *Response {
	return &Response{
		BaseResponse: &BaseResponse{
			Code: code,
			Msg:  msg,
		},
	}
}

func NewResponseFromPb(base *responsepb.BaseResponse) *Response {
	return &Response{
		BaseResponse: &BaseResponse{
			Code: int(base.Code),
			Msg:  base.Msg,
		},
	}
}

func NewResponseFromCode(code responsepb.Code) *Response {
	return &Response{
		BaseResponse: &BaseResponse{
			Code: int(code),
			Msg:  code.String(),
		},
	}
}

func (r *Response) SetData(data interface{}) *Response {
	r.Data = data

	return r
}

func (r *Response) SetMsg(msg string) *Response {
	r.Msg = msg
	return r
}
