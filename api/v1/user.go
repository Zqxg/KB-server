package v1

type CaptchaData struct {
	CaptchaId     string `json:"captchaId"`
	CaptchaBase64 string `json:"CaptchaBase64"`
	CaptchaAnswer string `json:"captchaAnswer"`
}
type CaptchaResponseData struct {
	CaptchaId     string `json:"captchaId"`
	CaptchaBase64 string `json:"CaptchaBase64"`
}

type RegisterRequest struct {
	Phone         string `json:"phone" binding:"required" example:"10012239028"`
	Password      string `json:"password" binding:"required" example:"123456"`
	CaptchaId     string `json:"captchaId" binding:"required"`     // 验证码ID字段
	CaptchaAnswer string `json:"captchaAnswer" binding:"required"` // 验证码字段
}

type PasswordLoginRequest struct {
	Phone         string `json:"phone" binding:"required" example:"10012239028"`
	Password      string `json:"password" binding:"required" example:"123456"`
	CaptchaId     string `json:"captchaId" binding:"required"`     // 验证码ID字段
	CaptchaAnswer string `json:"captchaAnswer" binding:"required"` // 验证码字段
}
type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" example:"alan"`
	Email    string `json:"email" binding:"email" example:"1234@gmail.com"`
}
type GetUserInfoResponseData struct {
	UserId    string `json:"userId"`
	Phone     string `json:"phone" example:"10012239028"`
	Nickname  string `json:"nickname" example:"alan"`
	RoleType  int    `json:"roleType" example:"0"`
	Email     string `json:"email"`
	CollegeId uint   `json:"collegeId"`
	StudentId string `json:"studentId"`
}
type UserAuthRequest struct {
	CollegeId uint   `json:"collegeId"`
	StudentId string `json:"studentId"`
	Remarks   string `json:"remarks"`
}
