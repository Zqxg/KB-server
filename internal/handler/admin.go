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
// @Router /v1/admin/createPublicKB [post]
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
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, v1.CreateKBResp{
		KBID: kbId,
	})
}

// UpdatePublicKB godoc
// @Summary 修改公共知识库
// @Schemes
// @Description
// @Tags 管理员模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UpdateKBNameReq true "params"
// @Success 200 {object} v1.Response
// @Router /v1/admin/updatePublicKB [post]
func (h *AdminHandler) UpdatePublicKB(ctx *gin.Context) {
	var req v1.UpdateKBNameReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	role := GetRoleTypeFromCtx(ctx)
	if role != enums.SUPER_ADMIN {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrPermissionDenied, nil)
	}
	err := h.kbService.UpdatePublicKB(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// DeletePublicKB godoc
// @Summary 删除公共知识库
// @Schemes
// @Description
// @Tags 管理员模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.DeleteKBReq true "params"
// @Success 200 {object} v1.Response
// @Router /v1/admin/deletePublicKB [post]
func (h *AdminHandler) DeletePublicKB(ctx *gin.Context) {
	var req v1.DeleteKBReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	role := GetRoleTypeFromCtx(ctx)
	if role != enums.SUPER_ADMIN {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrPermissionDenied, nil)
	}
	err := h.kbService.DeletePublicKB(ctx, req.KBID)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}
