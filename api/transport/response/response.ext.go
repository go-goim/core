// Code Written Manually

package response

import (
	"fmt"
)

/*
 * Define BaseResponse
 */

var _ error = &BaseResponse{}

func NewBaseResponse(code Code) *BaseResponse {
	return &BaseResponse{
		Code:   code,
		Reason: code.String(),
	}
}

func NewBaseResponseWithMessage(code Code, msg string) *BaseResponse {
	return &BaseResponse{
		Code:    code,
		Reason:  code.String(),
		Message: msg,
	}
}

func NewBaseResponseWithError(err error) *BaseResponse {
	return &BaseResponse{
		Code:    Code_InternalError,
		Reason:  Code_InternalError.String(),
		Message: err.Error(),
	}
}

func (x *BaseResponse) Error() string {
	return fmt.Sprintf("Code: %d, Reason: %s, Message: %s", x.Code, x.Reason, x.Message)
}

func (x *BaseResponse) Success() bool {
	return x.Code == Code_OK
}

func (x *BaseResponse) SetMessage(msg string) *BaseResponse {
	x.Message = msg
	return x
}

/*
 * Code
 */

func (x Code) BaseResponse() *BaseResponse {
	return NewBaseResponse(x)
}

func (x Code) BaseResponseWithMessage(msg string) *BaseResponse {
	return NewBaseResponseWithMessage(x, msg)
}

func (x Code) BaseResponseWithError(err error) *BaseResponse {
	return NewBaseResponseWithMessage(x, err.Error())
}
