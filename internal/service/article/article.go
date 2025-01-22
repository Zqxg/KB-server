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
	GetArticle(ctx context.Context, userId string, id uint) (*v1.ArticleData, error)
	CreateArticle(ctx context.Context, req *v1.CreateArticleRequest) (int, error)
	GetArticleCategory(ctx context.Context) ([]vo.CategoryView, error)
	UpdateArticle(ctx context.Context, req *v1.UpdateArticleRequest) (*v1.ArticleData, error)
	DeleteArticle(ctx context.Context, id uint) (int, error)
	DeleteArticleList(ctx context.Context, req *v1.DelArticleListReq) (int, error)
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

func (s *articleService) GetArticle(ctx context.Context, userId string, id uint) (*v1.ArticleData, error) {
	// 获取文章
	article, err := s.articleRepository.GetArticle(ctx, id)
	if err != nil {
		return nil, v1.ErrArticleNotExist
	}
	if article.Status != enums.StatusPublished && article.Status != enums.StatusPendingReview {
		return nil, v1.ErrArticleStatusError
	}
	// 私有 判断是否为本人
	if utils.Contains(article.VisibleRange, "private") {
		if userId != article.UserID {
			return nil, v1.ErrPermissionDenied
		}
	}
	Author, _ := s.userRepo.GetByUserId(ctx, article.UserID)
	category, _ := s.articleRepository.GetCategory(ctx, article.CategoryID)

	// 映射
	articleData := &v1.ArticleData{
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
		Status:          article.Status,
		CreatedAt:       utils.TimeFormat(article.CreatedAt, utils.FormatDateTime),
		UpdatedAt:       utils.TimeFormat(article.UpdatedAt, utils.FormatDateTime),
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

func (s *articleService) UpdateArticle(ctx context.Context, req *v1.UpdateArticleRequest) (*v1.ArticleData, error) {
	article, err := s.articleRepository.GetArticle(ctx, req.ArticleID)
	if err != nil {
		return nil, v1.ErrArticleNotExist
	}
	// 更新文章
	article.Title = req.Title
	article.Content = req.Content
	article.ContentShort = req.ContentShort
	article.CategoryID = req.CategoryID
	article.Importance = req.Importance
	article.VisibleRange = req.VisibleRange
	article.CommentDisabled = req.CommentDisabled
	article.SourceURI = req.SourceURI
	article.Status = enums.StatusPublished // todo：后续设置审核开关
	updateArticle, err := s.articleRepository.UpdateArticle(ctx, article)
	if err != nil {
		return nil, v1.ErrUpdateArticleFailed
	}
	// 映射
	Author, _ := s.userRepo.GetByUserId(ctx, article.UserID)
	category, _ := s.articleRepository.GetCategory(ctx, article.CategoryID)
	articleData := &v1.ArticleData{
		ArticleID:       updateArticle.ArticleID,
		Title:           updateArticle.Title,
		Content:         updateArticle.Content,
		ContentShort:    updateArticle.ContentShort,
		Author:          Author.Nickname,
		Category:        category.CategoryName,
		Importance:      updateArticle.Importance,
		VisibleRange:    updateArticle.VisibleRange,
		CommentDisabled: updateArticle.CommentDisabled,
		SourceURI:       updateArticle.SourceURI,
		UploadedFiles:   updateArticle.UploadedFiles,
		Status:          updateArticle.Status,
		CreatedAt:       utils.TimeFormat(updateArticle.CreatedAt, utils.FormatDateTime),
		UpdatedAt:       utils.TimeFormat(updateArticle.UpdatedAt, utils.FormatDateTime),
	}
	return articleData, nil
}

func (s *articleService) DeleteArticle(ctx context.Context, id uint) (int, error) {
	// 判断文章是否存在
	article, err := s.articleRepository.GetArticle(ctx, id)
	if err != nil {
		return -1, v1.ErrArticleNotExist
	}
	// 删除文章
	deletedCount, err := s.articleRepository.DeleteArticle(ctx, article.ArticleID)
	if err != nil {
		return -1, v1.ErrDeleteFailed
	}
	return deletedCount, nil
}

func (s *articleService) DeleteArticleList(ctx context.Context, req *v1.DelArticleListReq) (int, error) {
	// 批量删除文章
	deletedCount, err := s.articleRepository.DeleteArticleList(ctx, req.ArticleIDList)
	if err != nil {
		return -1, v1.ErrDeleteFailed
	}
	return deletedCount, nil
}
