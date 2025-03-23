package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"gorm.io/gorm"
	v1 "projectName/api/v1"
	"projectName/internal/enums"
	"projectName/internal/model"
	"time"
)

type ArticleRepository interface {
	GetArticle(ctx context.Context, id uint) (*model.Article, error)
	CreateArticle(ctx context.Context, article *model.Article) (int, error)
	GetArticleByTitleAndUserId(ctx context.Context, title string, authorID string) (*model.Article, error)
	UpdateArticle(ctx context.Context, article *model.Article) (*model.Article, error)
	DeleteArticle(ctx context.Context, id uint) (int, error)
	DeleteArticleList(ctx context.Context, ids []uint) (int, error)
	GetArticleListByCategory(ctx context.Context, categoryId uint, pageNum int, pageSize int) ([]model.Article, int64, error)
	GetUserArticleList(ctx context.Context, userId string, req *v1.GetUserArticleListReq, pageNum int, pageSize int) ([]model.Article, int64, error)
	GetArticleListByEs(ctx context.Context, indices []string, query *elastic.BoolQuery, highlight *elastic.Highlight, from, size int) (*elastic.SearchResult, error)
	CreateEsArticle(ctx context.Context, index string, article *model.EsArticle) error
	UpdateEsArticle(ctx context.Context, index string, article *model.EsArticle) error
	DeleteEsArticle(ctx context.Context, index string, articleId uint) error
	CreateEsIndex(ctx context.Context, index string) error // 新增es索引
	DeleteEsIndex(ctx context.Context, index string) error // 删除es索引
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
	r.logger.WithContext(ctx).Info("Successfully fetched knowledgeBase list", zap.Int("articleCount", len(articles)))

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
	if req.Status != -1 {
		query = query.Where("status =?", req.Status)
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

// GetArticleListByEs 根据多个 ES 索引查询文章
func (r *Repository) GetArticleListByEs(ctx context.Context, indices []string, query *elastic.BoolQuery, highlight *elastic.Highlight, from, size int) (*elastic.SearchResult, error) {
	r.logger.WithContext(ctx).Info("ES查询索引列表", zap.Any("indices", indices))

	searchService := r.esClient.Search().
		Index(indices...). // 支持多个索引
		Query(query).
		Highlight(highlight).
		From(from).Size(size) // 分页

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("ArticleRepository.GetArticleListByEs 查询失败", zap.Error(err))
		return nil, fmt.Errorf("failed to execute Elasticsearch query: %w", err)
	}

	r.logger.WithContext(ctx).Info("ArticleRepository.GetArticleListByEs 查询成功", zap.Any("total", searchResult.Hits.TotalHits.Value))
	return searchResult, nil
}

func (r *Repository) CreateEsArticle(ctx context.Context, index string, article *model.EsArticle) error {
	r.logger.WithContext(ctx).Info("ES index", zap.Any("index", index))
	_, err := r.esClient.Index().
		Index(index).
		Id(fmt.Sprintf("%d", article.ArticleID)).
		BodyJson(article).
		Do(ctx)
	r.logger.WithContext(ctx).Info("ArticleRepository.CreateEsArticle", zap.Any("knowledgeBase", article))
	if err != nil {
		r.logger.WithContext(ctx).Error("ArticleRepository.CreateEsArticle error", zap.Error(err))
		return fmt.Errorf("failed to create Elasticsearch document: %w", err)
	}
	return nil
}

func (r *Repository) UpdateEsArticle(ctx context.Context, index string, article *model.EsArticle) error {
	r.logger.WithContext(ctx).Info("ES index", zap.Any("index", index))
	_, err := r.esClient.Update().
		Index(index).
		Id(fmt.Sprintf("%d", article.ArticleID)).
		Doc(article).
		Do(ctx)
	r.logger.WithContext(ctx).Info("ArticleRepository.UpdateEsArticle", zap.Any("knowledgeBase", article))
	if err != nil {
		r.logger.WithContext(ctx).Error("ArticleRepository.UpdateEsArticle error", zap.Error(err))
		return fmt.Errorf("failed to update Elasticsearch document: %w", err)
	}
	return nil
}
func (r *Repository) DeleteEsArticle(ctx context.Context, index string, articleId uint) error {
	r.logger.WithContext(ctx).Info("ES index", zap.Any("index", index))
	_, err := r.esClient.Delete().
		Index(index).
		Id(fmt.Sprintf("%d", articleId)).
		Do(ctx)
	r.logger.WithContext(ctx).Info("ArticleRepository.DeleteEsArticle", zap.Any("articleId", articleId))
	if err != nil {
		r.logger.WithContext(ctx).Error("ArticleRepository.DeleteEsArticle error", zap.Error(err))
		return fmt.Errorf("failed to delete Elasticsearch document: %w", err)
	}
	return nil
}

func (r *Repository) CreateEsIndex(ctx context.Context, index string) error {
	// 创建索引
	createIndex, err := r.esClient.CreateIndex(index).Do(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("ArticleRepository.CreateEsIndex error", zap.Error(err))
		return fmt.Errorf("failed to create Elasticsearch index: %w", err)
	}
	if !createIndex.Acknowledged {
		r.logger.WithContext(ctx).Error("ArticleRepository.CreateEsIndex error", zap.Error(err))
		return fmt.Errorf("failed to create Elasticsearch index: %w", err)
	}
	return nil
}
func (r *Repository) DeleteEsIndex(ctx context.Context, index string) error {
	// 删除索引
	deleteIndex, err := r.esClient.DeleteIndex(index).Do(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("ArticleRepository.DeleteEsIndex error", zap.Error(err))
		return fmt.Errorf("failed to delete Elasticsearch index: %w", err)
	}
	if !deleteIndex.Acknowledged {
		r.logger.WithContext(ctx).Error("ArticleRepository.DeleteEsIndex error", zap.Error(err))
		return fmt.Errorf("failed to delete Elasticsearch index: %w", err)
	}
	return nil
}
