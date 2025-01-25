package repository

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	v1 "projectName/api/v1"
	"projectName/internal/enums"
	"projectName/internal/model"
	"projectName/internal/model/vo"
	"time"
)

type ArticleRepository interface {
	GetArticle(ctx context.Context, id uint) (*model.Article, error)
	CreateArticle(ctx context.Context, article *model.Article) (int, error)
	GetArticleByTitleAndUserId(ctx context.Context, title string, authorID string) (*model.Article, error)
	FetchAllCategoriesAndBuildTree(ctx context.Context) ([]vo.CategoryView, error)
	GetCategory(ctx context.Context, id uint) (*vo.CategoryView, error)
	UpdateArticle(ctx context.Context, article *model.Article) (*model.Article, error)
	DeleteArticle(ctx context.Context, id uint) (int, error)
	DeleteArticleList(ctx context.Context, ids []uint) (int, error)
	GetArticleListByCategory(ctx context.Context, categoryId uint, pageNum int, pageSize int) ([]model.Article, int64, error)
	GetUserArticleList(ctx context.Context, userId string, req *v1.GetUserArticleListReq, pageNum int, pageSize int) ([]model.Article, int64, error)
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

func (r *articleRepository) GetArticle(ctx context.Context, id uint) (*model.Article, error) {
	var article model.Article
	if err := r.DB(ctx).Table("kb_article").Where("article_id = ?", id).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrNotFound
		}
		r.logger.WithContext(ctx).Error("ArticleRepository.GetArticle error", zap.Error(err))
		return nil, err
	}

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

func (r *articleRepository) GetCategory(ctx context.Context, id uint) (*vo.CategoryView, error) {
	var categoryView vo.CategoryView

	// 查询视图中的单个分类，使用传入的 id
	if err := r.DB(ctx).Table("view_category_tree").Where("category_id = ?", id).First(&categoryView).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到记录，返回自定义的错误
			return nil, v1.ErrNotFound
		}
		r.logger.WithContext(ctx).Error("Failed to fetch category from view_category_tree", zap.Uint("categoryId", id), zap.Error(err))
		return nil, err
	}

	// 返回查询到的分类信息
	return &categoryView, nil
}

func (r *articleRepository) UpdateArticle(ctx context.Context, article *model.Article) (*model.Article, error) {
	if err := r.DB(ctx).Table("kb_article").Save(article).Error; err != nil {
		r.logger.WithContext(ctx).Error("ArticleRepository.UpdateArticle error", zap.Error(err))
		return nil, err
	}
	return article, nil
}

func (r *articleRepository) DeleteArticle(ctx context.Context, id uint) (int, error) {
	// 更新文章的 status 字段为已删除状态
	result := r.DB(ctx).Table("kb_article").
		Where("article_id = ?", id).
		Updates(map[string]interface{}{"status": enums.StatusDeleted})
	if result.Error != nil {
		r.logger.WithContext(ctx).Error("ArticleRepository.DeleteArticle error", zap.Error(result.Error))
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func (r *articleRepository) DeleteArticleList(ctx context.Context, ids []uint) (int, error) {
	// 更新文章的 status 字段为已删除状态
	updateResult := r.DB(ctx).Table("kb_article").
		Where("article_id IN (?)", ids).
		Update("status", enums.StatusDeleted)

	if updateResult.Error != nil {
		r.logger.WithContext(ctx).Error("ArticleRepository.DeleteArticleList UpdateStatus error", zap.Error(updateResult.Error))
		return 0, updateResult.Error
	}

	return int(updateResult.RowsAffected), nil
}

func (r *articleRepository) GetArticleListByCategory(ctx context.Context, categoryId uint, pageNum int, pageSize int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	// 计算偏移量
	offset := (pageNum - 1) * pageSize

	// 查询总数
	countResult := r.DB(ctx).Table("kb_article").
		Where("category_id = ? AND status = ?", categoryId, enums.StatusPublished).
		Count(&total)
	if countResult.Error != nil {
		r.logger.WithContext(ctx).Error("ArticleRepository.GetArticleListByCategory Count error", zap.Error(countResult.Error))
		return nil, 0, countResult.Error
	}

	// 如果没有数据，直接返回空
	if total == 0 {
		r.logger.WithContext(ctx).Info("No articles found", zap.Uint("categoryId", categoryId))
		return []model.Article{}, 0, nil
	}

	// 查询文章列表
	result := r.DB(ctx).Table("kb_article").
		Where("category_id = ? AND status = ?", categoryId, enums.StatusPublished).
		Offset(offset).
		Limit(pageSize).
		Find(&articles)
	if result.Error != nil {
		r.logger.WithContext(ctx).Error("ArticleRepository.GetArticleListByCategory Find error", zap.Error(result.Error))
		return nil, 0, result.Error
	}

	// 打印成功日志
	r.logger.WithContext(ctx).Info("Successfully fetched article list", zap.Int("articleCount", len(articles)))

	return articles, total, nil
}
func (r *articleRepository) GetUserArticleList(ctx context.Context, userId string, req *v1.GetUserArticleListReq, pageNum int, pageSize int) ([]model.Article, int64, error) {
	// 使用 GORM 获取数据库连接
	db := r.db.WithContext(ctx)

	// 创建查询构造器，开始构建查询条件
	query := db.Model(&model.Article{}).Where("user_id = ?", userId)

	// 根据请求参数添加查询条件
	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.CategoryID != 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}
	if req.CreatedAt != "" {
		// 转换字符串到时间类型并比较
		createdAt, err := time.Parse("2006-01-02", req.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid createdAt format: %v", err)
		}
		query = query.Where("created_at >= ?", createdAt)
	}
	if req.CreatedEnd != "" {
		// 转换字符串到时间类型并比较
		createdEnd, err := time.Parse("2006-01-02", req.CreatedEnd)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid CreatedEnd format: %v", err)
		}
		query = query.Where("created_at <= ?", createdEnd)
	}

	// 获取符合条件的文章总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 计算分页
	offset := (pageNum - 1) * pageSize

	// 获取分页后的文章列表
	var articles []model.Article
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	// 返回文章列表和总数
	return articles, total, nil
}
