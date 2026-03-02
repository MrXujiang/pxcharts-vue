package errorx

const (
	ErrNotFound         = 100001 // 查询时不存在
	ErrAlreadyExists    = 100002 // 查询时发现已存在
	ErrNoPermission     = 100003 // 无权限
	ErrInvalidParam     = 100004 // 请求参数错误
	ErrOperationFailed  = 100005 // 操作失败
	ErrOperationTimeout = 100006 // 操作超时

	// 用户相关 1001xx

	ErrRegister     = 100101 // 注册失败
	ErrInvalidEmail = 100102 // 邮箱格式不正确

	// 资源相关 1002xx

	ErrResourceNeedPassword  = 100201 // 资源需要密码
	ErrResourcePasswordError = 100202 // 资源密码错误
)
