package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	v1 "projectName/api/v1"
	"projectName/internal/enums"
	"projectName/internal/service/knowledgeBase"
	"projectName/pkg/utils"
)

type ArticleHandler struct {
	*Handler
	articleService knowledgeBase.ArticleService
}

func NewArticleHandler(
	handler *Handler,
	articleService knowledgeBase.ArticleService,
) *ArticleHandler {
	return &ArticleHandler{
		Handler:        handler,
		articleService: articleService,
	}
}

// CreateArticle godoc
// @Summary 新建文章
// @Schemes
// @Description
// @Tags 文章模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.CreateArticleRequest true "params"
// @Success 200 {object} v1.CreateArticleResponseData
// @Router /v1/article/createArticle [post]
func (h *ArticleHandler) CreateArticle(ctx *gin.Context) {
	var req v1.CreateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	articleId, err := h.articleService.CreateArticle(ctx, &req)
	if articleId == -1 || err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, v1.CreateArticleResponseData{
		ArticleID: articleId,
	})

}

// GetArticle godoc
// @Summary 获取文章详细
// @Schemes
// @Description
// @Tags 文章模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query int true "Article ID"
// @Success 200 {object} v1.ArticleData
// @Router /v1/article/getArticle [get]
func (h *ArticleHandler) GetArticle(ctx *gin.Context) {
	// 从查询参数中获取参数
	if !utils.IsNumeric(ctx.Query("id")) {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	articleID, _ := utils.ToInt(ctx.Query("id")) // 获取 articleID 参数
	if articleID < 0 {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	userId := GetUserIdFromCtx(ctx)
	articleData, err := h.articleService.GetArticle(ctx, userId, uint(articleID))
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, articleData)
}

// UpdateArticle godoc
// @Summary 修改文章内容
// @Schemes
// @Description
// @Tags 文章模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UpdateArticleRequest true "params"
// @Success 200 {object} v1.ArticleData
// @Router /v1/article/updateArticle [post]
func (h *ArticleHandler) UpdateArticle(ctx *gin.Context) {
	var req v1.UpdateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	userId, role := GetUserIdAndRoleTypeFromCtx(ctx)
	if req.AuthorID == userId || role == enums.SUPER_ADMIN {
		articleData, err := h.articleService.UpdateArticle(ctx, &req)
		if err != nil {
			v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
			return
		}
		v1.HandleSuccess(ctx, articleData)
	} else {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}
}

// DeleteArticleList godoc
// @Summary 批量删除文章
// @Schemes
// @Description
// @Tags 文章模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.DelArticleListReq true "params"
// @Success 200 {object} v1.DeleteArticleResponseData
// @Router /v1/article/deleteArticleList [post]
func (h *ArticleHandler) DeleteArticleList(ctx *gin.Context) {
	var req v1.DelArticleListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	role := GetRoleTypeFromCtx(ctx)
	if role == enums.SUPER_ADMIN {
		deletedCount, err := h.articleService.DeleteArticleList(ctx, &req)
		if err != nil {
			v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
			return
		}
		v1.HandleSuccess(ctx, v1.DeleteArticleResponseData{
			DeletedCount: deletedCount,
		})
	} else {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}
}

// DeleteArticle godoc
// @Summary 删除文章
// @Schemes
// @Description
// @Tags 文章模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.DeleteArticleRequest true "params"
// @Success 200 {object} v1.DeleteArticleResponseData
// @Router /v1/article/deleteArticle [post]
func (h *ArticleHandler) DeleteArticle(ctx *gin.Context) {
	var req v1.DeleteArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	userId, role := GetUserIdAndRoleTypeFromCtx(ctx)
	articleData, err := h.articleService.GetArticleById(ctx, req.ArticleID)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	if articleData.UserID == userId || role == enums.SUPER_ADMIN {
		deletedCount, err := h.articleService.DeleteArticle(ctx, req.ArticleID)
		if err != nil {
			v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
			return
		}
		v1.HandleSuccess(ctx, v1.DeleteArticleResponseData{
			DeletedCount: deletedCount,
		})
	} else {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrPermissionDenied, nil)
		return
	}
}

// GetUserArticleList godoc
// @Summary 获取个人文章列表
// @Schemes
// @Description
// @Tags 文章模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.GetUserArticleListReq true "params"
// @Success 200 {object} v1.ArticleList
// @Router /v1/article/getUserArticleList [post]
func (h *ArticleHandler) GetUserArticleList(ctx *gin.Context) {
	var req v1.GetUserArticleListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	userId := GetUserIdFromCtx(ctx)
	articleList, err := h.articleService.GetUserArticleList(ctx, userId, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, articleList)
}

// GetArticleListByEs godoc
// @Summary es文章查询
// @Schemes
// @Description
// @Tags 文章模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.GetArticleListByEsReq true "params"
// @Success 200 {object} v1.SearchArticleResp
// @Router /v1/article/getArticleListByEs [post]
func (h *ArticleHandler) GetArticleListByEs(ctx *gin.Context) {
	var req v1.GetArticleListByEsReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	// 默认值
	if req.Column == "" {
		req.Column = "_score"
	}
	if req.Order == "" {
		req.Order = "desc"
	}
	userId := GetUserIdFromCtx(ctx)

	articleList, err := h.articleService.GetArticleListByEs(ctx, userId, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
	}
	v1.HandleSuccess(ctx, articleList)
}

// GetArticleListByCID godoc
// @Summary 分类获取公开文章列表
// @Schemes
// @Description
// @Tags 文章模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.GetArticleListByCategoryReq true "params"
// @Success 200 {object} v1.ArticleList
// @Router /article/getArticleListByCID [post]
func (h *ArticleHandler) GetArticleListByCID(ctx *gin.Context) {
	var req v1.GetArticleListByCategoryReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	if req.CategoryID == 0 {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	userId, role := GetUserIdAndRoleTypeFromCtx(ctx)
	articleList, err := h.articleService.GetArticleListByCategory(ctx, userId, role, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, articleList)

}
