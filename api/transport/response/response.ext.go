// Code Written Manually

package response

import (
	"fmt"
)

/*
 * Define BaseResponse
 */

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

func (x *BaseResponse) SetCode(code Code) *BaseResponse {
	x.Code = code
	return x
}

func (x *BaseResponse) SetMsg(msg string) *BaseResponse {
	x.Msg = msg
	return x
}

/*
 * Define ResponseCode
 */

func (x Code) BaseResponse() *BaseResponse {
	return NewBaseResponse(x, x.String())
}
