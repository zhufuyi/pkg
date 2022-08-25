package errcode

import (
	"fmt"

	pb "github.com/zhufuyi/pkg/grpc/errcode/proto/rpcerrorpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error 错误
type Error struct {
	code int
	msg  string
}

var errorCodes = map[int]string{}

// NewError 创建新错误信息
func NewError(code int, msg string) *Error {
	if _, ok := errorCodes[code]; ok {
		panic(fmt.Sprintf("code %d 已经存在", code))
	}

	errorCodes[code] = msg

	return &Error{code: code, msg: msg}
}

// Code 错误码
func (e *Error) Code() int {
	return e.code
}

// Msg 错误信息
func (e *Error) Msg() string {
	return e.msg
}

// String 打印错误
func (e *Error) String() string {
	return fmt.Sprintf("code: %d, msg: %s", e.code, e.msg)
}

// ToRPCCode 自定义错误码转换为RPC识别的错误码，避免返回Unknown状态码
func ToRPCCode(code int) codes.Code {
	var statusCode codes.Code

	switch code {
	case InternalServerError.code:
		statusCode = codes.Internal
	case InvalidParam.code:
		statusCode = codes.InvalidArgument
	case Unauthorized.code:
		statusCode = codes.Unauthenticated
	case NotFound.code:
		statusCode = codes.NotFound
	case DeadlineExceeded.code:
		statusCode = codes.DeadlineExceeded
	case AccessDenied.code:
		statusCode = codes.PermissionDenied
	case LimitExceed.code:
		statusCode = codes.ResourceExhausted
	case MethodNotAllowed.code:
		statusCode = codes.Unimplemented
	default:
		statusCode = codes.Unknown
	}

	return statusCode
}

// ----------------------------------------------------------------------------------

// Status 状态
type Status struct {
	*status.Status
}

// FromError 把错误转换为status
func FromError(err error) *Status {
	s, _ := status.FromError(err)
	return &Status{s}
}

// ToGRPCStatus 除了原始业务错误码，新增其他说明信息msg，主要给内部客户端
func ToGRPCStatus(err *Error, msg string) *Status {
	s, _ := status.New(ToRPCCode(err.code), msg).WithDetails(&pb.Error{Code: int32(err.code), Message: err.msg})
	return &Status{s}
}

// ToGRPCError 通过Details属性返回错误信息给外部客户端
func ToGRPCError(err *Error) error {
	s, _ := status.New(ToRPCCode(err.code), err.msg).WithDetails(&pb.Error{Code: int32(err.code), Message: err.msg})
	return s.Err()
}
