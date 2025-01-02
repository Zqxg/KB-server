package v1

type RegisterRequest struct {
	Phone     string `json:"phone" binding:"required,phone" example:"10012239028"`
	Password  string `json:"password" binding:"required" example:"123456"`
	CaptchaId string `json:"captchaId"` // 验证码ID字段
	Captcha   string `json:"captcha"`   // 验证码字段
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required,phone" example:"10012239028"`
	Password string `json:"password" binding:"required" example:"123456"`
	Captcha  string `json:"captcha"` // 验证码字段
}
type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
}
type LoginResponse struct {
	Response
	Data LoginResponseData
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" example:"alan"`
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
}
type GetProfileResponseData struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname" example:"alan"`
}
type GetProfileResponse struct {
	Response
	Data GetProfileResponseData
}
