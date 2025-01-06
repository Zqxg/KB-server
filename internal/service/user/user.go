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
	"projectName/pkg/utils"
	"regexp"
	"strconv"
	"time"

	"github.com/DanPlayer/randomname"
)

type UserService interface {
	Register(ctx context.Context, req *v1.RegisterRequest) error
	PasswordLogin(ctx context.Context, req *v1.PasswordLoginRequest) (string, error)
	GetUserInfo(ctx context.Context, userId string) (*v1.GetUserInfoResponseData, error)
	UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error
	Logout(ctx context.Context, userId string, roleType int) error
	Cancel(ctx context.Context, userId string) error
	UserAuth(ctx context.Context, req *v1.UserAuthRequest, userId string, roleType int) error
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
	if !utils.IsPhoneNumber(req.Phone) {
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
		return v1.ErrDatabase
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
	if !utils.IsPhoneNumber(req.Phone) {
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
		key := fmt.Sprintf("%s%s:%d", enums.LOGIN_TOKEN_KEY, user.UserId, user.RoleType) // key="loginTokenKey:547519779070593342:0"
		if err = s.userRepo.Set(ctx, key, token, time.Hour*24*7); err != nil {
			return v1.ErrGetTokenFail // 存储失败
		}
		return nil
	})

	return token, nil
}

func (s *userService) GetUserInfo(ctx context.Context, userId string) (*v1.GetUserInfoResponseData, error) {
	user, err := s.userRepo.GetByUserId(ctx, userId)
	if err != nil {
		return nil, v1.ErrUserNotExist
	}

	return &v1.GetUserInfoResponseData{
		UserId:    user.UserId,
		Nickname:  user.Nickname,
		Phone:     user.Phone,
		RoleType:  user.RoleType,
		Email:     user.Email,
		CollegeId: user.CollegeId,
		StudentId: user.StudentId,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error {
	user, err := s.userRepo.GetByUserId(ctx, userId)
	if err != nil {
		return v1.ErrUserNotExist
	}
	if utils.IsEmpty(req.Email) && utils.IsEmpty(req.Nickname) {
		return v1.ErrParamEmpty
	}
	if utils.IsNotEmpty(req.Email) && utils.IsEmail(req.Email) {
		user.Email = req.Email
	}
	if utils.IsNotEmpty(req.Nickname) {
		user.Nickname = req.Nickname
	}
	if err = s.userRepo.Update(ctx, user); err != nil {
		return v1.ErrUpdateFailed
	}
	return nil
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
	key := fmt.Sprintf("%s%s:%d", enums.LOGIN_TOKEN_KEY, userId, roleType)
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

func (s *userService) UserAuth(ctx context.Context, req *v1.UserAuthRequest, userId string, roleType int) error {
	// 校验参数，确保是普通用户
	if roleType != enums.COMMON_USER {
		return v1.ErrUserAlreadyAuth
	}
	// 查询是否已有待处理认证请求
	existingAuthRequest, err := s.userRepo.GetUserAuthByUserId(ctx, userId)
	if err != nil {
		return v1.ErrDatabase
	}
	// 如果有待处理的认证请求，阻止重新发起认证
	if existingAuthRequest != nil && existingAuthRequest.Status == enums.WAITING {
		return v1.ErrUserAuthPending
	}
	// 认证请求
	userAuth := &model.UserAuth{
		UserId:      userId,
		RequestType: enums.SUTDENT_USER,
		Status:      enums.WAITING,
		ApplyTime:   time.Now(),
		CollegeId:   &req.CollegeId,
		StudentId:   &req.StudentId,
		Remarks:     &req.Remarks,
	}
	// 判断是否是第一次认证，或者认证请求状态是已拒绝或认证失败
	if existingAuthRequest == nil {
		// 第一次认证，存储认证请求信息到数据库
		if err = s.userRepo.CreateUserAuth(ctx, userAuth); err != nil {
			return v1.ErrUserAuthFailed
		}
	} else if existingAuthRequest.Status == enums.REJECTED || existingAuthRequest.Status == enums.FAILED {
		// 认证请求状态是已拒绝或认证失败，更新认证请求信息到数据库
		if err = s.userRepo.UpdateUserAuth(ctx, userAuth); err != nil {
			return v1.ErrUserAuthFailed
		}
	}
	return nil
}
