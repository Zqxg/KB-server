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

// GetArticle godoc
// @Summary 获取文章
// @Schemes
// @Description
// @Tags 文章模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.GetArticleRequest true "params"
// @Success 200 {object} v1.ArticleResponseData
func (h *ArticleHandler) CreateArticle(ctx *gin.Context) {
	var req v1.GetArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	if err := h.articleService.CreateArticle(ctx, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)

}
