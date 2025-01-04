package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"projectName/api/v1"
	"projectName/internal/service/user"
)

type UserHandler struct {
	*Handler
	captchaService user.CaptchaService
	userService    user.UserService
}

func NewUserHandler(handler *Handler, captchaService user.CaptchaService, userService user.UserService) *UserHandler {
	return &UserHandler{
		Handler:        handler,
		captchaService: captchaService,
		userService:    userService,
	}
}

// GetCaptcha godoc
// @Summary 获取验证码
// @Schemes
// @Description 获取验证码生成所需的ID和图片URL
// @Tags 用户模块
// @Accept json
// @Produce json
// @Success 200 {object} v1.CaptchaResponseData
// @Router /getCaptcha [get]
func (h *UserHandler) GetCaptcha(ctx *gin.Context) {
	captchaData, err := h.captchaService.GenerateCaptcha()
	if err != nil {
		h.logger.WithContext(ctx).Error("userService.GetCaptcha error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	// todo：不返回验证码答案 CaptchaResponseData
	v1.HandleSuccess(ctx, captchaData)
}

// Register godoc
// @Summary 用户注册
// @Schemes
// @Description 目前只支持邮箱登录
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.RegisterRequest true "params"
// @Success 200 {object} v1.Response
// @Router /register [post]
func (h *UserHandler) Register(ctx *gin.Context) {
	req := new(v1.RegisterRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.Register(ctx, req); err != nil {
		h.logger.WithContext(ctx).Error("userService.Register error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// PasswordLogin godoc
// @Summary 账号密码登录
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.PasswordLoginRequest true "params"
// @Success 200 {object} v1.LoginResponseData
// @Router /passwordLogin [post]
func (h *UserHandler) PasswordLogin(ctx *gin.Context) {
	var req v1.PasswordLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	token, err := h.userService.PasswordLogin(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	v1.HandleSuccess(ctx, v1.LoginResponseData{
		AccessToken: token,
	})
}

// GetProfile godoc
// @Summary 获取用户信息
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.GetProfileResponseData
// @Router /user/getProfile [get]
func (h *UserHandler) GetProfile(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	userData, err := h.userService.GetProfile(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}

	v1.HandleSuccess(ctx, v1.GetProfileResponseData{
		UserId:    userData.UserId,
		Phone:     userData.Phone,
		Nickname:  userData.Nickname,
		RoleType:  userData.RoleType,
		Email:     userData.Email,
		CollegeId: userData.CollegeId,
		StudentId: userData.StudentId,
	})
}

// UpdateProfile godoc
// @Summary 修改用户信息
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UpdateProfileRequest true "params"
// @Success 200 {object} v1.Response
// @Router /user/updateProfile [post]
func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)

	var req v1.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.UpdateProfile(ctx, userId, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// Logout godoc
// @Summary 退出用户
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.Response
// @Router /user/logout [get]
func (h *UserHandler) Logout(ctx *gin.Context) {
	userId, roleTpye := GetUserIdAndRoleTypeFromCtx(ctx)
	if err := h.userService.Logout(ctx, userId, roleTpye); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// Cancel godoc
// @Summary 注销用户
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.Response
// @Router /cancel [get]
func (h *UserHandler) Cancel(ctx *gin.Context) {
	userId, roleTpye := GetUserIdAndRoleTypeFromCtx(ctx)
	// 退出
	if err := h.userService.Logout(ctx, userId, roleTpye); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	// 注销
	if err := h.userService.Cancel(ctx, userId); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// UserAuth godoc
// @Summary 用户认证
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UserAuthRequest true "params"
// @Success 200 {object} v1.Response
// @Router /user/userAuth [post]
func (h *UserHandler) UserAuth(ctx *gin.Context) {
	userId, roleTpye := GetUserIdAndRoleTypeFromCtx(ctx)
	var req v1.UserAuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	if err := h.userService.UserAuth(ctx, &req, userId, roleTpye); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}
