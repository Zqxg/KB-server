package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	v1 "projectName/api/v1"
	"projectName/internal/enums"
	"projectName/internal/service/admin"
	"projectName/internal/service/knowledgeBase"
)

type AdminHandler struct {
	*Handler
	adminService admin.AdminService
	kbService    knowledgeBase.KnowledgeBaseService
}

func NewAdminHandler(
	handler *Handler,
	adminService admin.AdminService,
	kbService knowledgeBase.KnowledgeBaseService,
) *AdminHandler {
	return &AdminHandler{
		Handler:      handler,
		adminService: adminService,
		kbService:    kbService,
	}
}

// CreatePublicKB godoc
// @Summary 创建公共知识库
// @Schemes
// @Description
// @Tags 管理员模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.CreateKBRequest true "params"
// @Success 200 {object} v1.CreateKBResp
// @Router /admin/createPublicKB [post]
func (h *AdminHandler) CreatePublicKB(ctx *gin.Context) {
	var req v1.CreateKBRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	userId, role := GetUserIdAndRoleTypeFromCtx(ctx)
	if role != enums.SUPER_ADMIN {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrPermissionDenied, nil)
	}
	kbId, err := h.kbService.CreatePublicKB(ctx, userId, &req)
	if kbId == -1 || err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, v1.CreateKBResp{
		KBID: kbId,
	})
}
