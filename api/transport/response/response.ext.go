// Code Written Manually

package response

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
)

type IResponse interface {
	GetOrNewMeta() *Meta
	// SetMeta sets the meta information, but does not overwrite existing meta information when merging
	SetMeta(*Meta) IResponse
	SetData(interface{}) IResponse
	// SetBaseResponse sets the base response, but won't set the meta information.
	// Call SetMeta to set the meta information.
	SetBaseResponse(*BaseResponse) IResponse
	Marshall() ([]byte, error)
}

var protoMarshalOpt = protojson.MarshalOptions{EmitUnpopulated: true, UseEnumNumbers: true}

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
func (x *BaseResponse) SetData(data interface{}) IResponse {
	// do not set the data field
	return x
}

func (x *BaseResponse) SetBaseResponse(br *BaseResponse) IResponse {
	x.Code = br.Code
	x.Msg = br.Msg

	return x
}

func (x *BaseResponse) Marshall() ([]byte, error) {
	return protoMarshalOpt.Marshal(x)
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
type Response struct {
	*BaseResponse `json:",inline"`
	Data          interface{} `json:"data"`
}

func NewResponse(br *BaseResponse) *Response {
	return &Response{
		BaseResponse: br,
	}

}

func (x *Response) SetData(data interface{}) IResponse {
	x.Data = data

	return x
}

func (x *Response) Marshall() ([]byte, error) {
	return protoMarshalOpt.Marshal(x)
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
	ErrUnknown       = NewBaseResponse(CodeUnknownError, "UNKNOWN_ERROR")
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
