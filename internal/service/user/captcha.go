package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dchest/captcha"
)

type CaptchaService interface {
	GenerateCaptcha() (string, string, error)                                // 生成验证码
	VerifyCaptcha(ctx context.Context, captchaId, inputCaptcha string) error // 验证验证码
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

// GenerateCaptcha 生成验证码
func (s *SimpleCaptchaService) GenerateCaptcha() (captchaId string, captchaImageUrl string, err error) {
	// 生成一个新的验证码ID
	captchaId = captcha.New()

	// 生成验证码图片URL
	captchaImageUrl = fmt.Sprintf("/captcha/%s.png", captchaId)

	// 返回验证码ID和图片的URL
	return captchaId, captchaImageUrl, nil
}

// VerifyCaptcha 验证验证码
func (s *SimpleCaptchaService) VerifyCaptcha(ctx context.Context, captchaId, inputCaptcha string) error {
	// 验证用户输入的验证码是否正确
	if captcha.VerifyString(captchaId, inputCaptcha) {
		// 验证成功
		return nil
	}

	// 验证失败，返回错误
	return errors.New("invalid captcha")
}
