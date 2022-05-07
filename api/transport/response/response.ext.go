// Code Written Manually

package response

import (
	"encoding/json"
	"fmt"
	"strconv"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// NewMeta returns a new Meta object
func NewMeta() *Meta {
	return &Meta{}
}

func (x *Meta) SetRequestID(id string) *Meta {
	x.RequestId = id

	return x
}

func (x *Meta) SetTotal(total int) *Meta {
	var t = int32(total)
	x.Total = &t

	return x
}

func (x *Meta) SetPaging(page, size int) *Meta {
	if x.Pagination == nil {
		x.Pagination = &Pagination{}
	}

	x.Pagination.Page = int32(page)
	x.Pagination.PageSize = int32(size)

	return x
}

func (x *Meta) SetExtra(key, value string) *Meta {
	if x.Extra == nil {
		x.Extra = make(map[string]string)
	}

	x.Extra[key] = value

	return x
}

func (x *Meta) SetExtraInt(key string, value int) *Meta {
	if x.Extra == nil {
		x.Extra = make(map[string]string)
	}

	x.Extra[key] = strconv.Itoa(value)

	return x
}

func (x *Meta) SetExtraMap(m map[string]string) *Meta {
	if x.Extra == nil {
		x.Extra = make(map[string]string)
	}

	for k, v := range m {
		x.Extra[k] = v
	}

	return x
}

func (x *Meta) Merge(src *Meta) *Meta {
	if src == nil {
		return x
	}

	if src.RequestId != "" {
		x.RequestId = src.RequestId
	}

	if src.Total != nil {
		x.Total = src.Total
	}

	if src.Pagination != nil {
		x.Pagination = src.Pagination
	}

	if src.Extra != nil {
		x.SetExtraMap(src.Extra)
	}

	return x
}

type IResponse interface {
	GetOrNewMeta() *Meta
	// SetMeta sets the meta information, but does not overwrite existing meta information when merging
	SetMeta(*Meta) IResponse
	SetData(interface{}) (IResponse, error)
	// SetBaseResponse sets the base response, but won't set the meta information.
	// Call SetMeta to set the meta information.
	SetBaseResponse(*BaseResponse) IResponse
	Marshall() ([]byte, error)
}

/*
 * Define BaseResponse
 */

var _ IResponse = &BaseResponse{}
var _ error = &BaseResponse{}

func NewBaseResponse(code Code, msg string) *BaseResponse {
	return &BaseResponse{
		Code: code,
		Msg:  msg,
	}
}

func (x *BaseResponse) Error() string {
	return fmt.Sprintf("Code: %d, Msg: %s", x.Code, x.Msg)
}

func (x *BaseResponse) Success() bool {
	return x.Code == Code_OK
}

func (x *BaseResponse) GetOrNewMeta() *Meta {
	if x.GetMeta() == nil {
		x.Meta = NewMeta()
	}

	return x.Meta
}

func (x *BaseResponse) SetMeta(meta *Meta) IResponse {
	if meta == nil {
		return x
	}

	x.GetOrNewMeta().Merge(meta)
	return x
}

// SetData sets the data field of the response.
// Don't use this method in any circumstances.
func (x *BaseResponse) SetData(data interface{}) (IResponse, error) {
	return nil, fmt.Errorf("BaseResponse does not support SetData")
}

func (x *BaseResponse) SetBaseResponse(br *BaseResponse) IResponse {
	x.Code = br.Code
	x.Msg = br.Msg

	return x
}

func (x *BaseResponse) Marshall() ([]byte, error) {
	return json.Marshal(x)
}

func (x *BaseResponse) SetMsg(msg string) *BaseResponse {
	x.Msg = msg

	return x
}

/*
 * Define Response
 */

var _ IResponse = &Response{}

// Response is the response object for the HTTP transport.
// The difference between Response and PbResponse is that
// Response is the response object for non-pb messages of data.
// PbResponse is the response object for pb messages of data,
// it uses `protojson` to marshal the data instead of `json`.
type Response struct {
	*BaseResponse `json:",inline"`
	Meta          *Meta       `json:"meta"`
	Data          interface{} `json:"data"`
}

func NewResponse(br *BaseResponse) *Response {
	return &Response{
		BaseResponse: br,
	}

}

func (x *Response) SetData(data interface{}) (IResponse, error) {
	x.Data = data

	return x, nil
}

/*
 * PbResponse
 */

var _ IResponse = &PbResponse{}

func NewPbResponse(br *BaseResponse) *PbResponse {
	return &PbResponse{
		Code: br.Code,
		Msg:  br.Msg,
	}
}

func (x *PbResponse) GetOrNewMeta() *Meta {
	if x.GetMeta() == nil {
		x.Meta = NewMeta()
	}

	return x.Meta
}

func (x *PbResponse) SetMeta(meta *Meta) IResponse {
	if meta == nil {
		return x
	}

	x.GetOrNewMeta().Merge(meta)
	return x
}

// SetData sets the data field of the response.
// It will be marshalled into a protobuf message.
func (x *PbResponse) SetData(i interface{}) (IResponse, error) {
	if i == nil {
		return x, nil
	}

	message, ok := i.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("data must be a protobuf message")
	}

	data, err := anypb.New(message)
	if err != nil {
		return nil, err
	}

	x.Data = data

	return x, nil
}

func (x *PbResponse) SetBaseResponse(br *BaseResponse) IResponse {
	x.Code = br.Code
	x.Msg = br.Msg

	return x
}

func (x *PbResponse) Marshall() ([]byte, error) {
	o := protojson.MarshalOptions{EmitUnpopulated: true, UseEnumNumbers: true}
	return o.Marshal(x)
}

/*
 * Define ResponseCode
 */

const (
	// Code_OK       Code = 0

	// common error codes

	CodeUnknownError  Code = 10001
	CodeInvalidParams Code = 10002

	// user error codes

	CodeUserNotFound          Code = 20001
	CodeUserExists            Code = 20002
	CodeInvalidUserOrPassword Code = 20003

	// relation error codes

	CodeRelationNotFound            Code = 30001
	CodeInvalidUpdateRelationAction Code = 30002

	// push server error codes

	CodeUserNotOnline Code = 40001
)

var (
	// common error messages

	OK               = NewBaseResponse(Code_OK, "OK")
	ErrUnknown       = NewBaseResponse(CodeUnknownError, "unknown error")
	ErrInvalidParams = NewBaseResponse(CodeInvalidParams, "INVALID_PARAMS")

	// user error messages

	ErrUserNotFound          = NewBaseResponse(CodeUserNotFound, "USER_NOT_FOUND")
	ErrUserExist             = NewBaseResponse(CodeUserExists, "USER_EXIST")
	ErrInvalidUserOrPassword = NewBaseResponse(CodeInvalidUserOrPassword, "INVALID_USER_OR_PASSWORD")

	// relation error

	ErrRelationNotFound            = NewBaseResponse(CodeRelationNotFound, "RELATION_NOT_FOUND")
	ErrInvalidUpdateRelationAction = NewBaseResponse(CodeInvalidUpdateRelationAction, "INVALID_UPDATE_RELATION_ACTION")

	// push server error

	ErrUserNotOnline = NewBaseResponse(CodeUserNotOnline, "USER_NOT_ONLINE")
)
