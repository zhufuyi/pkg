package errcode

// 业务自定义错误
var (
	ErrorGetUserListFail = NewError(30000001, "获取用户数据失败")

	ErrorCreateMarketFail    = NewError(30000101, "新增仪表盘失败")
	ErrorCreateSubMarketFail = NewError(30000102, "新增监控模块失败")

	ErrorQueryPrometheusFail = NewError(30000201, "查询prometheus失败")

	ErrorUpdateFail = NewError(30000301, "更新失败")
)
