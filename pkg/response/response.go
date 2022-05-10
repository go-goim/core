package response

import (
	responsepb "github.com/yusank/goim/api/transport/response"
)

type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Meta *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Extra    map[string]string `json:"extra,omitempty"`
}

func BaseResponseFromPb(pb *responsepb.BaseResponse) *BaseResponse {
	return &BaseResponse{
		Code: int(pb.Code),
		Msg:  pb.Msg,
		Meta: MetaFromPbMeta(pb.Meta),
	}
}

func MetaFromPbMeta(pb *responsepb.Meta) *Meta {
	if pb == nil {
		return nil
	}

	return &Meta{
		Total:    int(pb.GetTotal()),
		Page:     int(pb.GetPage()),
		PageSize: int(pb.GetPageSize()),
		Extra:    pb.GetExtra(),
	}
}

type Response struct {
	*BaseResponse
	Data interface{} `json:"data,omitempty"`
}

func NewResponse(code int, msg string, data interface{}) *Response {
	return &Response{
		BaseResponse: &BaseResponse{
			Code: code,
			Msg:  msg,
		},
		Data: data,
	}
}

func NewResponseFromPb(base *responsepb.BaseResponse) *Response {
	return &Response{
		BaseResponse: BaseResponseFromPb(base),
	}
}

func (r *Response) SetData(data interface{}) *Response {
	r.Data = data

	return r
}
