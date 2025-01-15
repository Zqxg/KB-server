package repository

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"projectName/internal/model"
)

type ArticleRepository interface {
	GetArticle(ctx context.Context, id int64) (*model.Article, error)
	CreateArticle(ctx context.Context, article *model.Article) (uint, error)
	GetArticleByTitleAndUserId(ctx context.Context, title string, authorID string) (*model.Article, error)
}

func NewArticleRepository(
	repository *Repository,
) ArticleRepository {
	return &articleRepository{
		Repository: repository,
	}
}

type articleRepository struct {
	*Repository
}

func (r *articleRepository) GetArticle(ctx context.Context, id int64) (*model.Article, error) {
	var article model.Article

	return &article, nil
}

func (r *articleRepository) CreateArticle(ctx context.Context, article *model.Article) (uint, error) {
	if err := r.DB(ctx).Table("kb_article").Create(article).Error; err != nil {
		r.logger.WithContext(ctx).Error("ArticleRepository.CreateArticle error", zap.Error(err))
		return -1, err
	}
	return article.ArticleID, nil
}

func (r *articleRepository) GetArticleByTitleAndUserId(ctx context.Context, title string, authorID string) (*model.Article, error) {
	var article model.Article
	result := r.db.WithContext(ctx).
		Where("title = ? AND author_id = ?", title, authorID).
		First(&article)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 没有找到匹配的记录
		}
		r.logger.WithContext(ctx).Error("ArticleRepository.GetArticleByTitleAndUserId error", zap.Error(result.Error))
		return nil, result.Error // 其他错误
	}

	return &article, nil
}
