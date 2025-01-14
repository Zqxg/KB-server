package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"projectName/internal/model"
)

type ArticleRepository interface {
	GetArticle(ctx context.Context, id int64) (*model.Article, error)
	CreateArticle(ctx context.Context, article *model.Article) error
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

func (r *articleRepository) CreateArticle(ctx context.Context, article *model.Article) error {
	if err := r.DB(ctx).Table("kb_article").Create(article).Error; err != nil {
		return err
	}
	return nil
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
		return nil, result.Error // 其他错误
	}

	return &article, nil
}
