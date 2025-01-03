package handler

import (
	"github.com/gin-gonic/gin"
	"projectName/pkg/jwt"
	"projectName/pkg/log"
)

type Handler struct {
	logger *log.Logger
}

func NewHandler(
	logger *log.Logger,
) *Handler {
	return &Handler{
		logger: logger,
	}
}
func GetUserIdFromCtx(ctx *gin.Context) string {
	v, exists := ctx.Get("claims")
	if !exists {
		return ""
	}
	return v.(*jwt.MyCustomClaims).UserId
}

func GetRoleTypeFromCtx(ctx *gin.Context) int {
	v, exists := ctx.Get("claims")
	if !exists {
		return -1
	}
	return v.(*jwt.MyCustomClaims).RoleType
}

func GetUserIdAndRoleTypeFromCtx(ctx *gin.Context) (string, int) {
	v, exists := ctx.Get("claims")
	if !exists {
		return "", -1
	}
	return v.(*jwt.MyCustomClaims).UserId, v.(*jwt.MyCustomClaims).RoleType
}
