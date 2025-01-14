package article

import (
	"context"
	v1 "projectName/api/v1"
	"projectName/internal/enums"
	"projectName/internal/model"
	"projectName/internal/repository"
	"projectName/internal/service"
)

type ArticleService interface {
	GetArticle(ctx context.Context, id int64) (*model.Article, error)
	CreateArticle(ctx context.Context, req *v1.GetArticleRequest) error
}

func NewArticleService(
	service *service.Service,
	articleRepository repository.ArticleRepository,
) ArticleService {
	return &articleService{
		Service:           service,
		articleRepository: articleRepository,
	}
}

type articleService struct {
	*service.Service
	articleRepository repository.ArticleRepository
}

func (s *articleService) GetArticle(ctx context.Context, id int64) (*model.Article, error) {
	return s.articleRepository.GetArticle(ctx, id)
}

func (s *articleService) CreateArticle(ctx context.Context, req *v1.GetArticleRequest) error {
	// 判断是否有重复的文章标题&userId
	article, _ := s.articleRepository.GetArticleByTitleAndUserId(ctx, req.Title, req.AuthorID)
	if article != nil {
		return v1.ErrArticleAlreadyExist
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
	if err := s.articleRepository.CreateArticle(ctx, article); err != nil {
		return v1.ErrCreateArticleFailed
	}
	return nil
}
