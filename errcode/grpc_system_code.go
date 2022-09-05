package errcode

import (
	"github.com/golang/protobuf/proto"
)

// nolint
// rpc系统级别错误码，有Err后缀
var (
	SuccessErr         = NewRPCErr(0, "ok")
	InvalidParamsErr   = NewRPCErr(10001, "参数错误")
	UnauthorizedErr    = NewRPCErr(10002, "认证错误")
	InternalServerErr  = NewRPCErr(10003, "服务内部错误")
	NotFoundErr        = NewRPCErr(10004, "资源不存在")
	AlreadyExistsErr   = NewRPCErr(10005, "资源已存在")
	TimeoutErr         = NewRPCErr(10006, "超时")
	TooManyRequestsErr = NewRPCErr(10007, "请求过多")
	ForbiddenErr       = NewRPCErr(10008, "拒绝访问")
	LimitExceedErr     = NewRPCErr(10009, "访问限制")

	DeadlineExceededErr = NewRPCErr(10010, "已超过最后期限")
	AccessDeniedErr     = NewRPCErr(10011, "拒绝访问")
	MethodNotAllowedErr = NewRPCErr(10012, "不允许使用的方法")
)

// KV 键值对
type KV = map[string]interface{}

// RPCErr rpc error
func RPCErr(req proto.Message, err *GRPCStatus, details ...KV) error {
	var dts []proto.Message
	for _, detail := range details {
		dts = append(dts, NewDetails(detail))
	}
	return err.WithDetails(dts...).Status(req).Err()
}
