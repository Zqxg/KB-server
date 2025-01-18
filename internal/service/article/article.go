package article

import (
	"context"
	v1 "projectName/api/v1"
	"projectName/internal/enums"
	"projectName/internal/model"
	"projectName/internal/model/vo"
	"projectName/internal/repository"
	"projectName/internal/service"
	"projectName/pkg/utils"
)

type ArticleService interface {
	GetArticle(ctx context.Context, userId string, id int) (*v1.ArticleResponseData, error)
	CreateArticle(ctx context.Context, req *v1.CreateArticleRequest) (int, error)
	GetArticleCategory(ctx context.Context) ([]vo.CategoryView, error)
}

func NewArticleService(
	service *service.Service,
	articleRepository repository.ArticleRepository,
	userRepo repository.UserRepository,
) ArticleService {
	return &articleService{
		Service:           service,
		articleRepository: articleRepository,
		userRepo:          userRepo,
	}
}

type articleService struct {
	*service.Service
	articleRepository repository.ArticleRepository
	userRepo          repository.UserRepository
}

func (s *articleService) GetArticle(ctx context.Context, userId string, id int) (*v1.ArticleResponseData, error) {
	// 获取文章
	article, err := s.articleRepository.GetArticle(ctx, id)
	if err != nil {
		return nil, v1.ErrArticleNotExist
	}
	// 私有 判断是否为本人
	if !(utils.Contains(article.VisibleRange, "private") && userId == article.UserID) {
		return nil, v1.ErrPermissionDenied
	}
	Author, _ := s.userRepo.GetByUserId(ctx, article.UserID)
	category, _ := s.articleRepository.GetCategory(ctx, article.CategoryID)

	// 映射
	articleData := &v1.ArticleResponseData{
		ArticleID:       article.ArticleID,
		Title:           article.Title,
		Content:         article.Content,
		ContentShort:    article.ContentShort,
		Author:          Author.Nickname,
		Category:        category.CategoryName,
		Importance:      article.Importance,
		VisibleRange:    article.VisibleRange,
		CommentDisabled: article.CommentDisabled,
		SourceURI:       article.SourceURI,
		UploadedFiles:   article.UploadedFiles,
	}
	return articleData, nil
}

func (s *articleService) CreateArticle(ctx context.Context, req *v1.CreateArticleRequest) (int, error) {
	// 判断是否有重复的文章标题&userId
	article, _ := s.articleRepository.GetArticleByTitleAndUserId(ctx, req.Title, req.AuthorID)
	if article != nil {
		return -1, v1.ErrArticleAlreadyExist
	}
	article = &model.Article{
		Title:           req.Title,
		Content:         req.Content,
		ContentShort:    req.ContentShort,
		UserID:          req.AuthorID,
		CategoryID:      req.CategoryID,
		Importance:      req.Importance,
		VisibleRange:    req.VisibleRange,
		CommentDisabled: req.CommentDisabled,
		SourceURI:       req.SourceURI,
		Status:          enums.StatusPublished, // todo：后续设置审核开关
	}
	// 创建新文章
	articleId, err := s.articleRepository.CreateArticle(ctx, article)
	if err != nil {
		return -1, v1.ErrCreateArticleFailed
	}
	return articleId, nil
}

func (s *articleService) GetArticleCategory(ctx context.Context) ([]vo.CategoryView, error) {
	categories, err := s.articleRepository.FetchAllCategoriesAndBuildTree(ctx)
	if err != nil {
		return nil, v1.ErrQueryFailed
	}
	return categories, err
}
