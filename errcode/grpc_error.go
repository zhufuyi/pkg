package errcode

import (
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// GRPCStatus grpc 状态
type GRPCStatus struct {
	status  *status.Status
	details []proto.Message
}

// NewRPCErr 新建一个rpc状态对象
func NewRPCErr(code codes.Code, msg string) *GRPCStatus {
	return &GRPCStatus{
		status: status.New(code, msg),
	}
}

// Status 返回rpc状态
func (g *GRPCStatus) Status(details ...proto.Message) *status.Status {
	details = append(details, g.details...)
	st, err := g.status.WithDetails(details...)
	if err != nil {
		return g.status
	}
	return st
}

// WithDetails 附加详情信息
func (g *GRPCStatus) WithDetails(details ...proto.Message) *GRPCStatus {
	g.details = details
	return g
}

// NewDetails 创建detail
func NewDetails(details map[string]interface{}) proto.Message {
	detailStruct, err := structpb.NewStruct(details)
	if err != nil {
		return nil
	}
	return detailStruct
}

// ToRPCCode 转换为RPC识别的错误码，避免返回Unknown状态码
func ToRPCCode(code int) codes.Code {
	var statusCode codes.Code

	switch code {
	case InternalServerError.code:
		statusCode = codes.Internal
	case InvalidParams.code:
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
