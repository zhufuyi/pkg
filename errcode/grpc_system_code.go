package errcode

// nolint
// rpc系统级别错误码，有status前缀
var (
	StatusSuccess = NewGRPCStatus(0, "ok")

	StatusInvalidParams       = NewGRPCStatus(10001, "参数错误")
	StatusUnauthorized        = NewGRPCStatus(10002, "认证错误")
	StatusInternalServerError = NewGRPCStatus(10003, "服务内部错误")
	StatusNotFound            = NewGRPCStatus(10004, "资源不存在")
	StatusAlreadyExists       = NewGRPCStatus(10005, "资源已存在")
	StatusTimeout             = NewGRPCStatus(10006, "超时")
	StatusTooManyRequests     = NewGRPCStatus(10007, "请求过多")
	StatusForbidden           = NewGRPCStatus(10008, "拒绝访问")
	StatusLimitExceed         = NewGRPCStatus(10009, "访问限制")

	StatusDeadlineExceeded = NewGRPCStatus(10010, "已超过最后期限")
	StatusAccessDenied     = NewGRPCStatus(10011, "拒绝访问")
	StatusMethodNotAllowed = NewGRPCStatus(10012, "不允许使用的方法")
)
