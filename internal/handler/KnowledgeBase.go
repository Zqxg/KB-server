package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	v1 "projectName/api/v1"
	"projectName/internal/service/knowledgeBase"
	"projectName/pkg/utils"
)

type KnowledgeBaseHandler struct {
	*Handler
	knowledgeBaseService knowledgeBase.KnowledgeBaseService
}

func NewKnowledgeBaseHandler(
	handler *Handler,
	knowledgeBaseService knowledgeBase.KnowledgeBaseService,
) *KnowledgeBaseHandler {
	return &KnowledgeBaseHandler{
		Handler:              handler,
		knowledgeBaseService: knowledgeBaseService,
	}
}

// CreateKB godoc
// @Summary 新建团队知识库
// @Schemes
// @Description
// @Tags 知识库模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param req body v1.CreateKBRequest true "params"
// @Success 200 {object} v1.CreateKBResp
// @Router /v1/createKB [post]
func (h *KnowledgeBaseHandler) CreateKB(ctx *gin.Context) {
	var req v1.CreateKBRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	userId := GetUserIdFromCtx(ctx)
	kbId, err := h.knowledgeBaseService.CreateKnowledgeBase(ctx, userId, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, v1.CreateKBResp{
		KBID: kbId,
	})

}

// UpdateKBName godoc
// @Summary 更新知识库名称
// @Schemes
// @Description
// @Tags 知识库模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param req body v1.UpdateKBNameReq true "params"
// @Success 200 {object} v1.Response
// @Router /v1/updateKBName [post]
func (h *KnowledgeBaseHandler) UpdateKBName(ctx *gin.Context) {
	var req v1.UpdateKBNameReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	userId := GetUserIdFromCtx(ctx)
	err := h.knowledgeBaseService.UpdateKBName(ctx, userId, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// DeleteKB godoc
// @Summary 删除知识库
// @Schemes
// @Description
// @Tags 知识库模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param req body v1.DeleteKBReq true "params"
// @Success 200 {object} v1.Response
// @Router /v1/deleteKB [POST]
func (h *KnowledgeBaseHandler) DeleteKB(ctx *gin.Context) {
	var req v1.DeleteKBReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	userId := GetUserIdFromCtx(ctx)
	err := h.knowledgeBaseService.DeleteKB(ctx, userId, req.KBID)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// GetKBInfo godoc
// @Summary 获取知识库信息
// @Schemes
// @Description
// @Tags 知识库模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param kb_id query uint true "知识库ID"
// @Success 200 {object} v1.GetKBInfoResp
// @Router /v1/getKBInfo [GET]
func (h *KnowledgeBaseHandler) GetKBInfo(ctx *gin.Context) {
	// 从查询参数中获取参数
	if !utils.IsNumeric(ctx.Query("kb_id")) {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	kbId, _ := utils.ToUint(ctx.Query("kb_id")) // 获取 kb_id 参数
	if kbId < 0 {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	kbInfo, err := h.knowledgeBaseService.GetKBInfo(ctx, kbId)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, kbInfo)
}

// GetKBListByTeamId godoc
// @Summary 获取团队知识库列表
// @Schemes
// @Description
// @Tags 知识库模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param req body v1.GetKBListByTeamIdReq true "params"
// @Success 200 {object} v1.GetKBListByTeamIdResp
// @Router /v1/getKBListByTeamId [post]
func (h *KnowledgeBaseHandler) GetKBListByTeamId(ctx *gin.Context) {
	var req v1.GetKBListByTeamIdReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	kbList, err := h.knowledgeBaseService.GetKBListByTeamId(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, kbList)
}
