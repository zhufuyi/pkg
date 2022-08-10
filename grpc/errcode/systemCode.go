package errcode

// 系统级别错误码
var (
	Success             = NewError(0, "ok")
	InternalServerError = NewError(10000, "内部错误")
	InvalidParam        = NewError(10001, "无效参数")
	Unauthorized        = NewError(10002, "认证错误")
	NotFound            = NewError(10003, "没有找到")
	Unknown             = NewError(10004, "未知")
	DeadlineExceeded    = NewError(10005, "超出最后止期限")
	AccessDenied        = NewError(10006, "访问被拒绝")
	LimitExceed         = NewError(10007, "访问限制")
	MethodNotAllowed    = NewError(10008, "不支持该方法")

	SignParam          = NewError(10009, "无效签名")
	Validation         = NewError(10011, "验证失败")
	Database           = NewError(10012, "数据库错误")
	TooManyRequests    = NewError(10013, "请求太多")
	InvalidTransaction = NewError(10014, "无效传输")
	ServiceUnavailable = NewError(10015, "服务不可用")
)
