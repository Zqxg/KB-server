package user

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	v1 "projectName/api/v1"
	"projectName/internal/enums"
	"projectName/internal/model"
	"projectName/internal/repository"
	"projectName/internal/service"
	"time"

	"github.com/DanPlayer/randomname"
)

type UserService interface {
	Register(ctx context.Context, req *v1.RegisterRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (string, error)
	GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error)
	UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error
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
	// 校验验证码
	if err := s.captchaService.VerifyCaptcha(ctx, req.CaptchaId, req.Captcha); err != nil {
		return v1.ErrInvalidCaptcha // 如果验证码验证失败，返回错误
	}
	// 校验用户手机号是否已注册
	user, err := s.userRepo.GetByPhone(ctx, req.Phone)
	if err != nil {
		return v1.ErrInternalServerError
	}
	if err == nil && user != nil {
		return v1.ErrEmailAlreadyUse
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// 生成user_id
	userId, err := s.Sid.Gen25PrefixUID()
	if err != nil {
		return err
	}
	user = &model.User{
		UserId:   userId,
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

func (s *userService) Login(ctx context.Context, req *v1.LoginRequest) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return "", v1.ErrUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", err
	}
	token, err := s.Jwt.GenToken(user.UserId, time.Now().Add(time.Hour*24*7))
	if err != nil {
		return "", err
	}

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
