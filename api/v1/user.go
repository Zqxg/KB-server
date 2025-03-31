package v1

type CaptchaData struct {
	CaptchaId     string `json:"captcha_id"`
	CaptchaBase64 string `json:"Captcha_base64"`
	CaptchaAnswer string `json:"captcha_answer"`
}
type CaptchaResponseData struct {
	CaptchaId     string `json:"captcha_id"`
	CaptchaBase64 string `json:"Captcha_base64"`
}

type RegisterRequest struct {
	Phone         string `json:"phone" binding:"required" example:"10012239028"`
	Password      string `json:"password" binding:"required" example:"123456"`
	CaptchaId     string `json:"captcha_id" binding:"required"`     // 验证码ID字段
	CaptchaAnswer string `json:"captcha_answer" binding:"required"` // 验证码字段
}

type PasswordLoginRequest struct {
	Phone         string `json:"phone" binding:"required" example:"10012239028"`
	Password      string `json:"password" binding:"required" example:"123456"`
	CaptchaId     string `json:"captcha_id" binding:"required"`     // 验证码ID字段
	CaptchaAnswer string `json:"captcha_answer" binding:"required"` // 验证码字段
}
type LoginResponseData struct {
	AccessToken string `json:"access_token"`
}

type UpdateProfileRequest struct {
	Nickname string `json:"nick_name" example:"alan"`
	Email    string `json:"email" binding:"email" example:"1234@gmail.com"`
}
type GetUserInfoResponseData struct {
	UserId    string `json:"user_id"`
	Phone     string `json:"phone" example:"10012239028"`
	Nickname  string `json:"nick_name" example:"alan"`
	RoleType  int    `json:"role_type" example:"0"`
	Email     string `json:"email"`
	CollegeId uint   `json:"college_id"`
	StudentId string `json:"student_id"`
}
type UserAuthRequest struct {
	CollegeId uint   `json:"college_id"`
	StudentId string `json:"student_id"`
	Remarks   string `json:"remarks"`
}
