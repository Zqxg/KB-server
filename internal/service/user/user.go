package user

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	v1 "projectName/api/v1"
	"projectName/internal/enums"
	"projectName/internal/model"
	"projectName/internal/repository"
	"projectName/internal/service"
	"regexp"
	"strconv"
	"time"

	"github.com/DanPlayer/randomname"
)

type UserService interface {
	Register(ctx context.Context, req *v1.RegisterRequest) error
	PasswordLogin(ctx context.Context, req *v1.PasswordLoginRequest) (string, error)
	GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error)
	UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error
	Logout(ctx context.Context, userId string, roleType int) error
	Cancel(ctx context.Context, userId string) error
}

func NewUserService(
	service *service.Service,
	userRepo repository.UserRepository,
	captchaService CaptchaService, // 在构造函数中传入验证码服务
) UserService {
	return &userService{
		userRepo:       userRepo,
		captchaService: captchaService, // 注入验证码服务
		Service:        service,
	}
}

type userService struct {
	userRepo       repository.UserRepository
	captchaService CaptchaService // 新增验证码服务
	*service.Service
}

func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	// 校验手机号格式
	if !isValidPhone(req.Phone) {
		return v1.ErrPhoneFormat
	}
	// 校验密码格式
	if !isValidPassword(req.Password) {
		return v1.ErrPasswordFormat
	}
	// 校验验证码
	if !s.captchaService.VerifyCaptcha(req.CaptchaId, req.CaptchaAnswer) {
		return v1.ErrInvalidCaptcha // 如果验证码验证失败，返回错误
	}
	// 校验用户手机号是否已注册
	user, err := s.userRepo.GetByPhone(ctx, req.Phone)
	if err != nil {
		return v1.ErrInternalServerError
	}
	if err == nil && user != nil {
		return v1.ErrPhoneAlreadyUse
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// 生成user_id
	userId, err := s.Sid.GenSonyflakeID()
	if err != nil {
		return err
	}
	user = &model.User{
		UserId:   strconv.FormatInt(userId, 10),
		Phone:    req.Phone,
		Password: string(hashedPassword),
		Nickname: randomname.GenerateName(), //随机生成用户昵称
		RoleType: enums.COMMON_USER,         // 未认证前为普通用户
	}
	// Transaction demo
	err = s.Tm.Transaction(ctx, func(ctx context.Context) error {
		// Create a user
		if err = s.userRepo.Create(ctx, user); err != nil {
			return err
		}
		// TODO: other repo
		return nil
	})
	return err
}

func (s *userService) PasswordLogin(ctx context.Context, req *v1.PasswordLoginRequest) (string, error) {
	// 校验参数
	if !isValidPhone(req.Phone) {
		return "", v1.ErrPhoneFormat
	}
	if !s.captchaService.VerifyCaptcha(req.CaptchaId, req.CaptchaAnswer) {
		return "", v1.ErrInvalidCaptcha // 如果验证码验证失败，返回错误
	}
	user, err := s.userRepo.GetByPhone(ctx, req.Phone)
	if err != nil || user == nil {
		return "", v1.ErrUserNotExist
	}
	// 校验密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", v1.ErrDecryptPassword
	}
	// 生成token 有效期7天
	//todo: 后续可以改成后台可配置的天数
	token, err := s.Jwt.GenToken(user.UserId, user.RoleType, time.Now().Add(time.Hour*24*7))
	if err != nil {
		return "", v1.ErrGetTokenFail
	}
	// token存redis，后续注销和退出时候删除
	err = s.Tm.Transaction(ctx, func(ctx context.Context) error {
		// 将 token 存储到 Redis，设置过期时间为 7 天
		key := fmt.Sprintf("%s:%d", user.UserId, user.RoleType) // key="10012249028:0"
		if err = s.userRepo.Set(ctx, key, token, time.Hour*24*7); err != nil {
			return v1.ErrGetTokenFail // 存储失败
		}
		return nil
	})

	return token, nil
}

func (s *userService) GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &v1.GetProfileResponseData{
		UserId:   user.UserId,
		Nickname: user.Nickname,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	user.Email = req.Email
	user.Nickname = req.Nickname

	if err = s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

// 校验手机号
func isValidPhone(phone string) bool {
	// 中国手机号正则
	phoneRegex := `^1[0-9]\d{9}$`
	matched, err := regexp.MatchString(phoneRegex, phone)
	if err != nil {
		// 如果正则匹配出错，返回 false
		return false
	}
	return matched
}

// 校验密码长度
func isValidPassword(password string) bool {
	// 密码长度正则：长度必须在8到16之间
	passwordRegex := `^.{8,16}$`
	matched, err := regexp.MatchString(passwordRegex, password)
	if err != nil {
		// 如果正则匹配出错，返回 false
		return false
	}
	return matched
}

func (s *userService) Logout(ctx context.Context, userId string, roleType int) error {
	// 从 Redis 中删除 token
	key := fmt.Sprintf("%s:%d", userId, roleType)
	if err := s.userRepo.Delete(ctx, key); err != nil {
		s.Logger.Error("userService.Logout error", zap.Error(err))
		return v1.ErrLogoutFail
	}
	return nil
}

func (s *userService) Cancel(ctx context.Context, userId string) error {
	if err := s.userRepo.DeleteByUserId(ctx, userId); err != nil {
		s.Logger.Error("userService.Cancel error", zap.Error(err))
		return v1.ErrCancelFail
	}
	return nil
}
