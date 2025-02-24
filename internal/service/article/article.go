package article

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	v1 "projectName/api/v1"
	"projectName/internal/enums"
	"projectName/internal/model"
	"projectName/internal/model/vo"
	"projectName/internal/repository"
	"projectName/internal/service"
	"projectName/pkg/utils"
	"strings"
)

type ArticleService interface {
	GetArticleById(ctx context.Context, id uint) (*model.Article, error)
	GetArticle(ctx context.Context, userId string, id uint) (*v1.ArticleData, error)
	CreateArticle(ctx context.Context, req *v1.CreateArticleRequest) (int, error)
	GetArticleCategory(ctx context.Context) ([]vo.CategoryView, error)
	UpdateArticle(ctx context.Context, req *v1.UpdateArticleRequest) (*v1.ArticleData, error)
	DeleteArticle(ctx context.Context, id uint) (int, error)
	DeleteArticleList(ctx context.Context, req *v1.DelArticleListReq) (int, error)
	GetArticleListByCategory(ctx context.Context, req *v1.GetArticleListByCategoryReq) (*v1.ArticleList, error)
	GetUserArticleList(ctx context.Context, userId string, req *v1.GetUserArticleListReq) (*v1.ArticleList, error)
	GetArticleListByEs(ctx context.Context, req *v1.GetArticleListByEsReq) (*v1.SearchArticleResp, error)
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

func (s *articleService) GetArticleById(ctx context.Context, id uint) (*model.Article, error) {
	// 获取文章
	article, err := s.articleRepository.GetArticle(ctx, id)
	if err != nil {
		return nil, v1.ErrArticleNotExist
	}
	return article, nil
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

	// 反序列化上传的文件列表
	var uploadedFiles []v1.FileUpload
	if len(article.UploadedFiles) > 0 {
		err = json.Unmarshal(article.UploadedFiles, &uploadedFiles)
		if err != nil {
			return nil, v1.ErrDeserializeFileFailed
		}
	}

	// 映射
	articleData := &v1.ArticleData{
		ArticleID:       article.ArticleID,
		Title:           article.Title,
		Content:         article.Content,
		ContentShort:    article.ContentShort,
		Author:          Author.Nickname,
		Category:        category.CategoryName,
		CategoryID:      article.CategoryID,
		Importance:      article.Importance,
		VisibleRange:    article.VisibleRange,
		CommentDisabled: article.CommentDisabled,
		SourceURI:       article.SourceURI,
		UploadedFiles:   uploadedFiles,
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
	uploadedFilesData, err := json.Marshal(req.UploadedFiles)
	if err != nil {
		return -1, v1.ErrUploadFileFailed
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
		UploadedFiles:   uploadedFilesData,
		Status:          enums.StatusPublished, // todo：后续设置审核开关
	}
	// 创建新文章
	articleId, err := s.articleRepository.CreateArticle(ctx, article)
	// 判断是否公开，如果公开则创建es文档
	if strings.Contains(article.VisibleRange, "public") {
		esArticle := &model.EsArticle{
			ArticleID:    uint(articleId),
			Title:        article.Title,
			Content:      article.Content,
			CategoryID:   article.CategoryID,
			UserID:       article.UserID,
			Status:       article.Status,
			VisibleRange: article.VisibleRange,
			CreatedAt:    article.CreatedAt,
			UpdatedAt:    article.UpdatedAt,
		}
		esArticle.ArticleID = uint(articleId)

		// 创建es文档
		if err = s.articleRepository.CreateEsArticle(ctx, esArticle); err != nil {
			return -1, v1.ErrCreateEsArticleFailed
		}
	}
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
	// 反序列化上传的文件列表
	var uploadedFiles []v1.FileUpload
	if len(article.UploadedFiles) > 0 {
		err = json.Unmarshal(article.UploadedFiles, &uploadedFiles)
		if err != nil {
			return nil, v1.ErrDeserializeFileFailed
		}
	}
	articleData := &v1.ArticleData{
		ArticleID:       updateArticle.ArticleID,
		Title:           updateArticle.Title,
		Content:         updateArticle.Content,
		ContentShort:    updateArticle.ContentShort,
		Author:          Author.Nickname,
		Category:        category.CategoryName,
		CategoryID:      updateArticle.CategoryID,
		Importance:      updateArticle.Importance,
		VisibleRange:    updateArticle.VisibleRange,
		CommentDisabled: updateArticle.CommentDisabled,
		SourceURI:       updateArticle.SourceURI,
		UploadedFiles:   uploadedFiles,
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

func (s *articleService) GetArticleListByCategory(ctx context.Context, req *v1.GetArticleListByCategoryReq) (*v1.ArticleList, error) {
	// 查询文章列表及分页信息
	pageIndex, pageSize := initPage(req.PageIndex, req.PageSize)
	articles, total, err := s.articleRepository.GetArticleListByCategory(ctx, req.CategoryID, pageIndex, pageSize)
	if err != nil {
		return nil, v1.ErrQueryFailed
	}

	// 映射文章数据
	var articleList []*v1.ArticleData
	for _, article := range articles {
		// 获取作者昵称
		Author, _ := s.userRepo.GetByUserId(ctx, article.UserID)
		// 获取分类名称
		category, _ := s.articleRepository.GetCategory(ctx, article.CategoryID)
		// 反序列化上传的文件列表
		var uploadedFiles []v1.FileUpload
		if len(article.UploadedFiles) > 0 {
			err = json.Unmarshal(article.UploadedFiles, &uploadedFiles)
			if err != nil {
				return nil, v1.ErrDeserializeFileFailed
			}
		}
		articleData := &v1.ArticleData{
			ArticleID:       article.ArticleID,
			Title:           article.Title,
			Content:         article.Content,
			ContentShort:    article.ContentShort,
			Author:          Author.Nickname,
			Category:        category.CategoryName,
			CategoryID:      article.CategoryID,
			Importance:      article.Importance,
			VisibleRange:    article.VisibleRange,
			CommentDisabled: article.CommentDisabled,
			SourceURI:       article.SourceURI,
			UploadedFiles:   uploadedFiles,
			Status:          article.Status,
			CreatedAt:       utils.TimeFormat(article.CreatedAt, utils.FormatDateTime),
			UpdatedAt:       utils.TimeFormat(article.UpdatedAt, utils.FormatDateTime),
		}
		articleList = append(articleList, articleData)
	}

	// 构建返回结构
	response := &v1.ArticleList{
		ArticleDataList: articleList,
		PageResponse: v1.PageResponse{
			TotalCount: total,
			PageIndex:  pageIndex,
			PageSize:   pageSize,
		},
	}

	return response, nil
}

// page初始化
func initPage(pageIndex int, pageSize int) (int, int) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 10 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return pageIndex, pageSize
}
func (s *articleService) GetUserArticleList(ctx context.Context, userId string, req *v1.GetUserArticleListReq) (*v1.ArticleList, error) {
	// 查询文章列表及分页信息
	pageIndex, pageSize := initPage(req.PageIndex, req.PageSize)
	// 查询文章列表
	articles, total, err := s.articleRepository.GetUserArticleList(ctx, userId, req, pageIndex, pageSize)
	if err != nil {
		return nil, v1.ErrQueryFailed
	}
	// 映射文章数据
	Author, _ := s.userRepo.GetByUserId(ctx, userId)
	var articleList []*v1.ArticleData
	for _, article := range articles {
		// 获取分类名称
		category, _ := s.articleRepository.GetCategory(ctx, article.CategoryID)
		// 反序列化上传的文件列表
		var uploadedFiles []v1.FileUpload
		if len(article.UploadedFiles) > 0 {
			err = json.Unmarshal(article.UploadedFiles, &uploadedFiles)
			if err != nil {
				return nil, v1.ErrDeserializeFileFailed
			}
		}
		articleData := &v1.ArticleData{
			ArticleID:       article.ArticleID,
			Title:           article.Title,
			Content:         article.Content,
			ContentShort:    article.ContentShort,
			Author:          Author.Nickname,
			Category:        category.CategoryName,
			CategoryID:      article.CategoryID,
			Importance:      article.Importance,
			VisibleRange:    article.VisibleRange,
			CommentDisabled: article.CommentDisabled,
			SourceURI:       article.SourceURI,
			UploadedFiles:   uploadedFiles,
			Status:          article.Status,
			CreatedAt:       utils.TimeFormat(article.CreatedAt, utils.FormatDateTime),
			UpdatedAt:       utils.TimeFormat(article.UpdatedAt, utils.FormatDateTime),
		}
		articleList = append(articleList, articleData)
	}
	// 构建返回结构
	response := &v1.ArticleList{
		ArticleDataList: articleList,
		PageResponse: v1.PageResponse{
			TotalCount: total,
			PageIndex:  pageIndex,
			PageSize:   pageSize,
		},
	}
	return response, nil
}

func (s *articleService) GetArticleListByEs(ctx context.Context, req *v1.GetArticleListByEsReq) (*v1.SearchArticleResp, error) {
	// 1. 设置分页信息
	pageNo, pageSize := initPage(req.PageIndex, req.PageSize)

	// 2. 构建查询条件
	query := elastic.NewBoolQuery()

	// 根据标题进行搜索 逗号分隔
	if req.Title != "" {
		query = query.Must(elastic.NewMatchQuery("title", req.Title))
	}

	// 根据内容进行搜索 逗号分隔
	if req.Content != "" {
		query = query.Must(elastic.NewMatchQuery("content", req.Content))
	}

	// 可以根据需要增加更多的过滤条件，例如状态（status）、重要性（importance）等
	if req.SearchMode == "1" {
		// todo: 高级搜索处理（例如，可能加上更多的过滤条件或权重设置）
	}

	// 根据关键字进行全文搜索 逗号分隔
	if len(req.Keywords) > 0 {
		for _, keyword := range req.Keywords {
			if req.PhraseMatch { // 启用短语匹配
				query = query.Should(elastic.NewMatchPhraseQuery("content", keyword))
			} else {
				query = query.Should(elastic.NewMatchQuery("content", keyword))
			}
		}
	}

	// 将 req.Categories 转换为 []interface{}
	var categories []interface{}
	for _, category := range req.Categories {
		categories = append(categories, category)
	}
	// 根据分类 ID 进行过滤（如果有）
	if len(req.Categories) > 0 {
		query = query.Filter(elastic.NewTermsQuery("category_id", categories...))
	}

	// 3. 添加高亮查询
	highlight := elastic.NewHighlight().Field("content").PreTags("<em>").PostTags("</em>")

	// 4. 设置分页查询
	from := (pageNo - 1) * pageSize

	// 5. 调用 repository 中的查询方法
	searchResult, err := s.articleRepository.GetArticleListByEs(ctx, query, highlight, from, pageSize)
	if err != nil {
		return nil, err
	}

	// 6. 解析搜索结果，构建响应数据
	var articles []v1.ArticleSearchInfo
	for _, hit := range searchResult.Hits.Hits {
		var esArticle model.EsArticle
		var article v1.ArticleSearchInfo
		if err = json.Unmarshal(hit.Source, &esArticle); err != nil {
			continue
		}
		// 获取并设置评分
		article.EsArticle = esArticle
		article.Score = *hit.Score

		// 获取高亮内容（如果有）
		if highlightFields, ok := hit.Highlight["content"]; ok {
			// 将高亮部分替换为 HTML 格式
			article.EsArticle.Content = strings.Join(highlightFields, "...")
		}

		articles = append(articles, article)
	}

	// 7. 构建分页响应
	resp := &v1.SearchArticleResp{
		PageResponse: v1.PageResponse{
			TotalCount: searchResult.Hits.TotalHits.Value,
			PageIndex:  pageNo,
			PageSize:   pageSize,
		},
		Articles: articles,
	}

	return resp, nil
}
