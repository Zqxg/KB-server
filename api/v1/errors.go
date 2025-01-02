package v1

var (
	// common errors
	ErrSuccess             = newError(0, "ok")
	ErrBadRequest          = newError(400, "Bad Request")
	ErrUnauthorized        = newError(401, "Unauthorized")
	ErrNotFound            = newError(404, "Not Found")
	ErrInternalServerError = newError(500, "Internal Server Error")

	// more biz errors
	ErrEmailAlreadyUse = newError(1001, "邮箱已存在")
	ErrPoneAlreadyUse  = newError(1002, "手机号已存在")

	// 2000 错误码
	ErrInvalidCaptcha = newError(2000, "验证码错误")

	// 3000 数据库
	ErrDatabase = newError(3000, "数据库错误")

	// 20000 业务逻辑错误
)
