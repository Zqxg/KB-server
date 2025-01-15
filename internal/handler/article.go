package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	v1 "projectName/api/v1"
	"projectName/internal/service/article"
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
