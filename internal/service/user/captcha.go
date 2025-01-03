package user

import (
	"github.com/mojocn/base64Captcha"
	v1 "projectName/api/v1"
	"time"
)

type CaptchaService interface {
	GenerateCaptcha() (v1.CaptchaData, error)                  // 生成验证码
	VerifyCaptcha(captchaId string, captchaAnswer string) bool // 验证验证码
}

type SimpleCaptchaService struct {
	captchaExpireDuration time.Duration // 验证码有效期
}

// NewCaptchaService 创建一个验证码服务实例，返回 CaptchaService 接口类型
func NewCaptchaService(expireDuration time.Duration) CaptchaService {
	return &SimpleCaptchaService{
		captchaExpireDuration: expireDuration,
	}
}

// 数字驱动生成验证码的配置
var digitDriver = base64Captcha.DriverDigit{
	Height:   50,
	Width:    200,
	Length:   4,   //验证码长度
	MaxSkew:  0.7, //倾斜
	DotCount: 1,   //背景的点数，越大，字体越模糊
}

// 设置自带的store
var store = base64Captcha.DefaultMemStore

// GenerateCaptcha 生成验证码
func (s *SimpleCaptchaService) GenerateCaptcha() (v1.CaptchaData, error) {
	var Data v1.CaptchaData
	// 创建一个新的验证码对象
	captcha := base64Captcha.NewCaptcha(&digitDriver, store)

	// 生成验证码ID、base64图片和答案
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		return Data, err
	}
	Data.CaptchaId = id
	Data.CaptchaBase64 = b64s
	Data.CaptchaAnswer = answer
	return Data, nil
}

// VerifyCaptcha 验证验证码
func (s *SimpleCaptchaService) VerifyCaptcha(captchaId string, captchaAnswer string) bool {
	return store.Verify(captchaId, captchaAnswer, true)
}
