package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	v1 "projectName/api/v1"
	"projectName/internal/enums"
	"projectName/internal/service/article"
	"projectName/pkg/utils"
)

type ArticleHandler struct {
	*Handler
	articleService article.ArticleService
}

func NewArticleHandler(
	handler *Handler,
	articleService article.ArticleService,
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
// @Router /article/createArticle [post]
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

// GetArticleCategory godoc
// @Summary 获取文章分组
// @Schemes
// @Description
// @Tags 文章模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.CategoryData
// @Router /article/getArticleCategory [get]
func (h *ArticleHandler) GetArticleCategory(ctx *gin.Context) {
	categories, err := h.articleService.GetArticleCategory(ctx)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, v1.CategoryData{
		CategoryList: categories,
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
// @Router /article/getArticle [get]
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
// @Router /article/UpdateArticle [post]
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
// @Router /article/DeleteArticleList [post]
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
// @Router /article/DeleteArticle [post]
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

// GetArticleListByCategory godoc
// @Summary 分类获取公开文章列表
// @Schemes
// @Description
// @Tags 文章模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param categoryId query int true "Category ID"
// @Param pageIndex query int true "Page Index"
// @Param pageSize query int true "Page Size"
// @Success 200 {object} v1.ArticleList
// @Router /article/getArticleListByCategory [get]
func (h *ArticleHandler) GetArticleListByCategory(ctx *gin.Context) {
	// 从查询参数中获取参数
	categoryId, _ := utils.ToInt(ctx.DefaultQuery("categoryId", "1")) // 获取 categoryId 参数
	pageIndex, _ := utils.ToInt(ctx.DefaultQuery("pageIndex", "1"))   // 获取 pageIndex 参数
	pageSize, _ := utils.ToInt(ctx.DefaultQuery("pageSize", "10"))    // 获取 pageSize 参数
	req := v1.GetArticleListByCategoryReq{
		CategoryID: uint(categoryId),
		PageRequest: v1.PageRequest{
			PageIndex: pageIndex,
			PageSize:  pageSize,
		},
	}
	articleList, err := h.articleService.GetArticleListByCategory(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, articleList)

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
// @Router /article/getUserArticleList [post]
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
