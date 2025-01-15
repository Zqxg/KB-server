package repository

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"projectName/internal/model"
	"projectName/internal/model/vo"
)

type ArticleRepository interface {
	GetArticle(ctx context.Context, id int64) (*model.Article, error)
	CreateArticle(ctx context.Context, article *model.Article) (int, error)
	GetArticleByTitleAndUserId(ctx context.Context, title string, authorID string) (*model.Article, error)
	FetchAllCategoriesAndBuildTree(ctx context.Context) ([]vo.CategoryView, error)
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

func (r *articleRepository) CreateArticle(ctx context.Context, article *model.Article) (int, error) {
	if err := r.DB(ctx).Table("kb_article").Create(article).Error; err != nil {
		r.logger.WithContext(ctx).Error("ArticleRepository.CreateArticle error", zap.Error(err))
		return -1, err
	}
	return int(article.ArticleID), nil
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

// BuildCategoryTree 用于将平坦的分类数据转换为树状结构
func BuildCategoryTree(data []vo.CategoryView, parentId uint) []vo.CategoryView {
	var result []vo.CategoryView
	for _, item := range data {
		if item.ParentId == parentId {
			item.Children = BuildCategoryTree(data, item.CId)
			result = append(result, item)
		}
	}
	return result
}

// FetchAllCategoriesAndBuildTree 从数据库获取所有分类数据并构建树状结构
func (r *articleRepository) FetchAllCategoriesAndBuildTree(ctx context.Context) ([]vo.CategoryView, error) {
	var categories []vo.CategoryView
	// 查询视图中的所有分类
	if err := r.DB(ctx).Table("view_category_tree").Find(&categories).Error; err != nil {
		r.logger.WithContext(ctx).Error("Failed to fetch categories from view_category_tree", zap.Error(err))
		return nil, err
	}
	// 调用 BuildCategoryTree 函数将平坦的分类数据转换为树状结构
	tree := BuildCategoryTree(categories, 0)
	r.logger.WithContext(ctx).Info("Successfully built category tree", zap.Int("rootCount", len(tree)))
	return tree, nil
}
