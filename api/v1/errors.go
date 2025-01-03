package v1

var (
	// common errors
	ErrSuccess             = newError(0, "操作成功")
	ErrBadRequest          = newError(400, "请求错误")
	ErrUnauthorized        = newError(401, "未授权")
	ErrPermissionDenied    = newError(403, "权限不足")
	ErrNotFound            = newError(404, "未找到")
	ErrInternalServerError = newError(500, "内部服务器错误")

	// more biz errors
	ErrEmailAlreadyUse = newError(1001, "邮箱已存在")
	ErrPhoneAlreadyUse = newError(1002, "手机号已存在")
	ErrPhoneFormat     = newError(1003, "手机号格式错误")
	ErrPasswordFormat  = newError(1004, "密码格式错误")
	ErrDecryptPassword = newError(1005, "错误解密密码")
	ErrGetTokenFail    = newError(1006, "获取token失败")
	ErrUserNotExist    = newError(1007, "用户不存在")
	ErrLogoutFail      = newError(1008, "用户退出失败")
	ErrCancelFail      = newError(1009, "用户注销失败")

	// 2000 错误码
	ErrInvalidCaptcha = newError(2000, "验证码错误")

	// 3000 数据库
	ErrDatabase = newError(3000, "数据库错误")

	// 20000 业务逻辑错误
	ErrParamEmpty = newError(20000, "参数为空")
)
