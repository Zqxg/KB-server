package knowledgeBase

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	v1 "projectName/api/v1"
	"projectName/internal/enums"
	"projectName/internal/model"
	"projectName/internal/repository"
	"projectName/internal/service"
	"projectName/pkg/utils"
	"strconv"
	"strings"
)

type ArticleService interface {
	GetArticleById(ctx context.Context, id uint) (*model.Article, error)
	GetArticle(ctx context.Context, userId string, id uint) (*v1.ArticleData, error)
	CreateArticle(ctx context.Context, role int, req *v1.CreateArticleRequest) (int, error)
	UpdateArticle(ctx context.Context, userID string, req *v1.UpdateArticleRequest) (*v1.ArticleData, error)
	DeleteArticle(ctx context.Context, userId string, role int, id uint) (int, error)
	DeleteArticleList(ctx context.Context, req *v1.DelArticleListReq) (int, error)
	GetArticleListByCategory(ctx context.Context, userID string, role int, req *v1.GetArticleListByCategoryReq) (*v1.ArticleList, error)
	GetUserArticleList(ctx context.Context, userId string, req *v1.GetUserArticleListReq) (*v1.ArticleList, error)
	GetArticleListByEs(ctx context.Context, userId string, req *v1.GetArticleListByEsReq) (*v1.SearchArticleResp, error)
}

func NewArticleService(
	service *service.Service,
	articleRepository repository.ArticleRepository,
	userRepo repository.UserRepository,
	kbRepository repository.KBRepository,
	teamRepository repository.TeamRepository,
) ArticleService {
	return &articleService{
		Service:           service,
		articleRepository: articleRepository,
		userRepo:          userRepo,
		kbRepository:      kbRepository,
		teamRepository:    teamRepository,
	}
}

type articleService struct {
	*service.Service
	articleRepository repository.ArticleRepository
	userRepo          repository.UserRepository
	kbRepository      repository.KBRepository
	teamRepository    repository.TeamRepository
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
	if article.Status == enums.StatusDeleted {
		return nil, v1.ErrArticleStatusError
	}
	// 获取知识库类型
	kb, _ := s.kbRepository.GetKBViewById(ctx, article.KBID)
	// 私人知识库，私人可见
	if kb.KBType == enums.KBTypePrivate {
		if userId != article.UserID {
			return nil, v1.ErrPermissionDenied
		}
	}
	// 团队知识库，团队成员可见
	if kb.KBType == enums.KBTypeTeam {
		// 获取当前团队成员列表
		teamMembers, err := s.teamRepository.GetTeamMemberListByTeamID(ctx, kb.TeamID)
		if err != nil {
			return nil, v1.ErrTeamNotExist
		}
		// 检查用户是否为团队成员
		var userIds []string
		for _, member := range teamMembers {
			userIds = append(userIds, member.UserID)
		}
		// 判断当前角色是否在团队成员列表中
		if !utils.ContainsString(userIds, userId) {
			return nil, v1.ErrPermissionDenied
		}
	}
	Author, _ := s.userRepo.GetByUserId(ctx, article.UserID)
	category, _ := s.kbRepository.GetCategoryById(ctx, article.CategoryID)

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
		TeamID:          kb.TeamID,
		TeamName:        kb.TeamName,
		KBID:            article.KBID,
		KBType:          kb.KBType,
		KBName:          kb.KbName,
		Importance:      article.Importance,
		CommentDisabled: article.CommentDisabled,
		SourceURI:       article.SourceURI,
		UploadedFiles:   uploadedFiles,
		Status:          article.Status,
		CreatedAt:       utils.TimeFormat(article.CreatedAt, utils.FormatDateTime),
		UpdatedAt:       utils.TimeFormat(article.UpdatedAt, utils.FormatDateTime),
	}
	return articleData, nil
}

func (s *articleService) CreateArticle(ctx context.Context, role int, req *v1.CreateArticleRequest) (int, error) {
	// 判断是否有重复的文章标题&userId
	article, _ := s.articleRepository.GetArticleByTitleAndUserId(ctx, req.Title, req.AuthorID)
	if article != nil {
		return -1, v1.ErrArticleAlreadyExist
	}
	// 判断知识库是否存在
	kb, _ := s.kbRepository.GetKBViewById(ctx, req.KBID)
	if kb == nil {
		return -1, v1.ErrKnowledgeNotExist
	}
	// 私人知识库，只有作者可以创建文章
	if kb.KBType == enums.KBTypePrivate {
		if kb.UserID != req.AuthorID {
			return -1, v1.ErrPermissionDenied
		}
	}
	if kb.KBType == enums.KBTypeTeam {
		// 团队知识库，只有团队成员可以创建文章
		// 获取当前团队成员列表
		teamMembers, err := s.teamRepository.GetTeamMemberListByTeamID(ctx, kb.TeamID)
		if err != nil {
			return -1, v1.ErrTeamNotExist
		}
		// 检查用户是否为团队成员
		var userIds []string
		for _, member := range teamMembers {
			userIds = append(userIds, member.UserID)
		}
		// 判断当前角色是否在团队成员列表中
		if !utils.ContainsString(userIds, req.AuthorID) {
			return -1, v1.ErrPermissionDenied
		}
	}
	if kb.KBType == enums.KBTypePublic && role != enums.SUPER_ADMIN {
		return -1, v1.ErrPermissionDenied
	}

	// 判断分类是否存在
	category, _ := s.kbRepository.GetCategoryById(ctx, req.CategoryID)
	if category == nil {
		return -1, v1.ErrCategoryNotExist
	}
	// 分类ID与知识库ID不匹配
	if category.KbID != req.KBID {
		return -1, v1.ErrCategoryNotMatchKB
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
		KBID:            req.KBID,
		CategoryID:      req.CategoryID,
		Importance:      req.Importance,
		CommentDisabled: req.CommentDisabled,
		SourceURI:       req.SourceURI,
		UploadedFiles:   uploadedFilesData,
		Status:          req.Status,
		UpdatedBy:       req.AuthorID,
	}
	// 创建新文章
	articleId, err := s.articleRepository.CreateArticle(ctx, article)
	if err != nil {
		return -1, v1.ErrCreateArticleFailed
	}
	// 判断文章状态
	if article.Status == enums.StatusPublished {
		// 判断知识库类型，选择es索引
		esIndex := s.GetESIndex(ctx, article.UserID, req.KBID)
		if esIndex == "" {
			return -1, v1.ErrCreateEsIndexFailed
		}
		// 创建es文档
		esArticle := &model.EsArticle{
			ArticleID:       uint(articleId),
			Title:           article.Title,
			Content:         article.Content,
			ContentShort:    article.ContentShort,
			KBID:            article.KBID,
			CategoryID:      article.CategoryID,
			UserID:          article.UserID,
			Importance:      article.Importance,
			CommentDisabled: article.CommentDisabled,
			SourceURI:       article.SourceURI,
			Status:          article.Status,
			UploadedFile:    false,
			CreatedAt:       article.CreatedAt,
			UpdatedAt:       article.UpdatedAt,
		}
		esArticle.ArticleID = uint(articleId)
		if article.UploadedFiles != nil {
			esArticle.UploadedFile = true
		}
		// 创建es文档
		if err = s.articleRepository.CreateEsArticle(ctx, esIndex, esArticle); err != nil {
			return -1, v1.ErrCreateEsArticleFailed
		}
	}
	return articleId, nil
}

func (s *articleService) UpdateArticle(ctx context.Context, userID string, req *v1.UpdateArticleRequest) (*v1.ArticleData, error) {
	// 查询旧文章（获取原 KBID）
	oldArticle, err := s.articleRepository.GetArticle(ctx, req.ArticleID)
	if err != nil {
		return nil, v1.ErrArticleNotExist
	}
	oldKBID := oldArticle.KBID
	oldStatus := oldArticle.Status

	// 更新文章内容
	oldArticle.Title = req.Title
	oldArticle.Content = req.Content
	oldArticle.ContentShort = req.ContentShort
	oldArticle.CategoryID = req.CategoryID
	oldArticle.Importance = req.Importance
	oldArticle.KBID = req.KBID // ⚠️ 这里可能改了 KBID
	oldArticle.CommentDisabled = req.CommentDisabled
	oldArticle.SourceURI = req.SourceURI
	oldArticle.Status = req.Status
	oldArticle.UpdatedBy = userID

	// 更新数据库
	updatedArticle, err := s.articleRepository.UpdateArticle(ctx, oldArticle)
	if err != nil {
		return nil, v1.ErrUpdateArticleFailed
	}

	// 查询作者信息 & 分类
	author, _ := s.userRepo.GetByUserId(ctx, updatedArticle.UserID)
	updatedAuthor, _ := s.userRepo.GetByUserId(ctx, updatedArticle.UpdatedBy)
	category, _ := s.kbRepository.GetCategoryById(ctx, updatedArticle.CategoryID)
	kb, _ := s.kbRepository.GetKBViewById(ctx, updatedArticle.KBID)

	// 处理上传文件
	var uploadedFiles []v1.FileUpload
	if len(updatedArticle.UploadedFiles) > 0 {
		if err = json.Unmarshal(updatedArticle.UploadedFiles, &uploadedFiles); err != nil {
			return nil, v1.ErrDeserializeFileFailed
		}
	}

	// 构造返回
	articleData := &v1.ArticleData{
		ArticleID:       updatedArticle.ArticleID,
		Title:           updatedArticle.Title,
		Content:         updatedArticle.Content,
		ContentShort:    updatedArticle.ContentShort,
		Author:          author.Nickname,
		Category:        category.CategoryName,
		CategoryID:      updatedArticle.CategoryID,
		TeamID:          kb.TeamID,
		TeamName:        kb.TeamName,
		KBType:          kb.KBType,
		Importance:      updatedArticle.Importance,
		KBID:            updatedArticle.KBID,
		KBName:          kb.KbName,
		CommentDisabled: updatedArticle.CommentDisabled,
		SourceURI:       updatedArticle.SourceURI,
		UploadedFiles:   uploadedFiles,
		Status:          updatedArticle.Status,
		UpdatedBy:       updatedAuthor.Nickname,
		CreatedAt:       utils.TimeFormat(updatedArticle.CreatedAt, utils.FormatDateTime),
		UpdatedAt:       utils.TimeFormat(updatedArticle.UpdatedAt, utils.FormatDateTime),
	}

	// 生成旧索引和新索引
	oldIndex := s.GetESIndex(ctx, updatedArticle.UserID, oldKBID)
	newIndex := s.GetESIndex(ctx, updatedArticle.UserID, updatedArticle.KBID)

	// 构建ES文章对象
	esArticle := &model.EsArticle{
		ArticleID:       updatedArticle.ArticleID,
		Title:           updatedArticle.Title,
		Content:         updatedArticle.Content,
		ContentShort:    updatedArticle.ContentShort,
		KBID:            updatedArticle.KBID,
		Importance:      updatedArticle.Importance,
		CommentDisabled: updatedArticle.CommentDisabled,
		SourceURI:       updatedArticle.SourceURI,
		CategoryID:      updatedArticle.CategoryID,
		UserID:          updatedArticle.UserID,
		Status:          updatedArticle.Status,
		UploadedFile:    updatedArticle.UploadedFiles != nil,
		CreatedAt:       updatedArticle.CreatedAt,
		UpdatedAt:       updatedArticle.UpdatedAt,
	}

	// 如果KBID变了，处理ES索引迁移
	if oldIndex != newIndex {
		// 删除旧索引下的文档
		_ = s.articleRepository.DeleteEsArticle(ctx, oldIndex, updatedArticle.ArticleID)
		// 新增到新索引
		if err = s.articleRepository.CreateEsArticle(ctx, newIndex, esArticle); err != nil {
			return nil, v1.ErrUpdateEsArticleFailed
		}
	} else {
		// KB未变,判断status
		if oldStatus == enums.StatusDraft && updatedArticle.Status == enums.StatusPublished {
			// 新增到新索引
			if err = s.articleRepository.CreateEsArticle(ctx, newIndex, esArticle); err != nil {
				return nil, v1.ErrUpdateEsArticleFailed
			}
		}
		if err = s.articleRepository.UpdateEsArticle(ctx, newIndex, esArticle); err != nil {
			return nil, v1.ErrUpdateEsArticleFailed
		}
	}

	return articleData, nil
}

func (s *articleService) DeleteArticle(ctx context.Context, userId string, role int, id uint) (int, error) {
	// 判断文章是否存在
	article, err := s.articleRepository.GetArticle(ctx, id)
	if err != nil {
		return -1, v1.ErrArticleNotExist
	}
	if role != enums.SUPER_ADMIN && userId != article.UserID {
		return -1, v1.ErrPermissionDenied
	}
	// 获取文章知识库所属团队
	kb, _ := s.kbRepository.GetKBViewById(ctx, article.KBID)
	// 判断当前用户是否为团队管理员或文章作者
	if kb.KBType == enums.KBTypeTeam {
		// 获取当前团队成员列表
		teamMembers, err := s.teamRepository.GetTeamMemberListByTeamID(ctx, kb.TeamID)
		if err != nil {
			return -1, v1.ErrTeamNotExist
		}
		// 检查用户是否为团队成员
		var userIds []string
		for _, member := range teamMembers {
			userIds = append(userIds, member.UserID)
		}
		// 判断当前用户是否在团队成员列表中
		if !utils.ContainsString(userIds, userId) {
			return -1, v1.ErrPermissionDenied
		}
	}
	// 判断文章状态
	if article.Status == enums.StatusPublished {
		// 判断知识库类型，选择es索引
		kb, _ := s.kbRepository.GetKBViewById(ctx, article.KBID)
		index := ""
		if kb.KBType == enums.KBTypePrivate {
			index = enums.Private_knowledge_index + article.UserID
		}
		if kb.KBType == enums.KBTypePublic {
			index = enums.Public_knowledge_index + strconv.Itoa(int(article.KBID))
		}
		if kb.KBType == enums.KBTypeTeam {
			index = enums.Team_knowledge_index + strconv.Itoa(int(kb.TeamID))
		}
		// 删除es文档
		if err = s.articleRepository.DeleteEsArticle(ctx, index, article.ArticleID); err != nil {
			return -1, v1.ErrDeleteEsArticleFailed
		}
	}
	// 删除文章
	deletedCount, err := s.articleRepository.DeleteArticle(ctx, article.ArticleID)
	if err != nil {
		return -1, v1.ErrDeleteFailed
	}

	return deletedCount, nil
}

// 暂时废弃
func (s *articleService) DeleteArticleList(ctx context.Context, req *v1.DelArticleListReq) (int, error) {
	// 批量删除文章
	deletedCount, err := s.articleRepository.DeleteArticleList(ctx, req.ArticleIDList)
	if err != nil {
		return -1, v1.ErrDeleteFailed
	}
	return deletedCount, nil
}

func (s *articleService) GetArticleListByCategory(ctx context.Context, userID string, role int, req *v1.GetArticleListByCategoryReq) (*v1.ArticleList, error) {
	// 查询文章列表及分页信息
	pageIndex, pageSize := service.InitPage(req.PageIndex, req.PageSize)
	// 判断分类是否存在
	category, err := s.kbRepository.GetCategoryById(ctx, req.CategoryID)
	if err != nil {
		return nil, v1.ErrCategoryNotExist
	}
	// 获取知识库类型
	kbView, err := s.kbRepository.GetKBViewById(ctx, category.KbID)
	if err != nil {
		return nil, v1.ErrKnowledgeNotExist
	}
	// 超级管理员可以查看所有文章
	if role != enums.SUPER_ADMIN {
		if kbView.KBType == enums.KBTypeTeam {
			// 校验用户是否在团队中
			_, err := s.teamRepository.GetMemberByTeamIDAndUserID(ctx, kbView.TeamID, userID)
			if err != nil {
				return nil, v1.ErrPermissionDenied
			}
		}
	}
	status := enums.GetStatus(req.Status)
	articles, total, err := s.articleRepository.GetArticleListByCategory(ctx, req.CategoryID, status, pageIndex, pageSize)
	if err != nil {
		return nil, v1.ErrQueryFailed
	}

	// 映射文章数据
	var articleList []*v1.ArticleData
	for _, article := range articles {
		// 获取作者昵称
		Author, _ := s.userRepo.GetByUserId(ctx, article.UserID)
		updatedBy, _ := s.userRepo.GetByUserId(ctx, article.UpdatedBy)
		// 获取分类名称
		category, _ := s.kbRepository.GetCategoryById(ctx, article.CategoryID)
		kb, _ := s.kbRepository.GetKBViewById(ctx, article.KBID)
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
			TeamID:          kb.TeamID,
			TeamName:        kb.TeamName,
			KBType:          kb.KBType,
			Importance:      article.Importance,
			KBID:            article.KBID,
			KBName:          kb.KbName,
			CommentDisabled: article.CommentDisabled,
			SourceURI:       article.SourceURI,
			UploadedFiles:   uploadedFiles,
			Status:          article.Status,
			UpdatedBy:       updatedBy.Nickname,
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

func (s *articleService) GetUserArticleList(ctx context.Context, userId string, req *v1.GetUserArticleListReq) (*v1.ArticleList, error) {
	// 查询文章列表及分页信息
	pageIndex, pageSize := service.InitPage(req.PageIndex, req.PageSize)
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
		category, _ := s.kbRepository.GetCategoryById(ctx, article.CategoryID)
		kb, _ := s.kbRepository.GetKBViewById(ctx, article.KBID)
		updatedBy, _ := s.userRepo.GetByUserId(ctx, article.UpdatedBy)
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
			TeamID:          kb.TeamID,
			TeamName:        kb.TeamName,
			KBType:          kb.KBType,
			KBID:            article.KBID,
			KBName:          kb.KbName,
			Importance:      article.Importance,
			CommentDisabled: article.CommentDisabled,
			SourceURI:       article.SourceURI,
			UploadedFiles:   uploadedFiles,
			Status:          article.Status,
			UpdatedBy:       updatedBy.Nickname,
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

func (s *articleService) GetArticleListByEs(ctx context.Context, userId string, req *v1.GetArticleListByEsReq) (*v1.SearchArticleResp, error) {
	// 1. 设置分页信息
	pageNo, pageSize := service.InitPage(req.PageIndex, req.PageSize)

	// 2. 构建查询条件
	query := elastic.NewBoolQuery()

	if req.AdvSearch { // 高级搜索，必须满足所有条件
		if req.Title != "" {
			query = query.Should(elastic.NewMatchQuery("title", req.Title))
		}
		if req.Content != "" {
			query = query.Should(elastic.NewMatchQuery("content", req.Content))
		}
		if len(req.Keywords) > 0 {
			for _, keyword := range req.Keywords {
				if req.PhraseMatch {
					query = query.Must(elastic.NewMatchPhraseQuery("content_short", keyword))
				} else {
					query = query.Should(elastic.NewMatchQuery("content_short", keyword))
				}
			}
		}
		if req.CreateTimeStart != "" && req.CreateTimeEnd != "" {
			query = query.Filter(elastic.NewRangeQuery("created_at").Gte(req.CreateTimeStart).Lte(req.CreateTimeEnd))
		}
		if importance, _ := utils.ToInt(req.Importance); importance > 0 {
			query = query.Filter(elastic.NewTermsQuery("importance", importance))
		}
		if req.KBID > 0 {
			query = query.Filter(elastic.NewTermsQuery("kb_id", req.KBID))
		}
	} else { // 普通搜索
		if len(req.Keywords) > 0 {
			for _, keyword := range req.Keywords {
				if req.PhraseMatch {
					query = query.Must(elastic.NewMatchPhraseQuery("content", keyword)).
						Must(elastic.NewMatchPhraseQuery("title", keyword)).
						Must(elastic.NewMatchPhraseQuery("content_short", keyword))
				} else {
					query = query.Should(elastic.NewMatchQuery("content", keyword)).
						Should(elastic.NewMatchQuery("title", keyword)).
						Should(elastic.NewMatchQuery("content_short", keyword))
				}
			}
		}
	}
	// 分类过滤
	if len(req.Categories) > 0 {
		var categories []interface{}
		for _, category := range req.Categories {
			categories = append(categories, category)
		}
		query = query.Filter(elastic.NewTermsQuery("category_id", categories...))
	}

	// 3. 添加高亮
	highlight := elastic.NewHighlight().
		Field("content").PreTags("<mark>").PostTags("</mark>").
		Field("title").PreTags("<mark>").PostTags("</mark>").
		Field("content_short").PreTags("<mark>").PostTags("</mark>")

	// 4. 计算分页
	from := (pageNo - 1) * pageSize

	// 5. 动态组装用户可查询的 ES 索引
	var indices []string
	// 公共知识库
	indices = append(indices, enums.Public_knowledge_index+"*")
	// 私人知识库
	indices = append(indices, enums.Private_knowledge_index+userId)
	// 团队知识库（需要查用户所在的团队ids）
	teamIds, _ := s.teamRepository.GetTeamIdsByUserID(ctx, userId)
	for _, teamId := range teamIds {
		indices = append(indices, enums.Team_knowledge_index+teamId)
	}

	// 6. 调用 ES 查询
	searchResult, err := s.articleRepository.GetArticleListByEs(ctx, indices, query, highlight, from, pageSize)
	if err != nil {
		return nil, err
	}

	// 7. 处理查询结果
	var articles []v1.ArticleSearchInfo
	for _, hit := range searchResult.Hits.Hits {
		var esArticle model.EsArticle
		if err := json.Unmarshal(hit.Source, &esArticle); err != nil {
			continue
		}
		kb, _ := s.kbRepository.GetKBViewById(ctx, esArticle.KBID)
		user, _ := s.userRepo.GetByUserId(ctx, esArticle.UserID)
		category, _ := s.kbRepository.GetCategoryById(ctx, esArticle.CategoryID)

		article := v1.ArticleSearchInfo{
			Title:           esArticle.Title,
			Content:         esArticle.Content,
			ContentShort:    esArticle.ContentShort,
			UploadedFile:    esArticle.UploadedFile,
			Status:          esArticle.Status,
			Author:          user.Nickname,
			KBName:          kb.KbName,
			TeamName:        kb.TeamName,
			Category:        category.CategoryName,
			ArticleID:       esArticle.ArticleID,
			CreatedAt:       esArticle.CreatedAt,
			UpdatedAt:       esArticle.UpdatedAt,
			Importance:      esArticle.Importance,
			CommentDisabled: esArticle.CommentDisabled,
			SourceURI:       esArticle.SourceURI,
		}
		article.Score = *hit.Score

		if highlightFields, ok := hit.Highlight["content"]; ok {
			article.Content = strings.Join(highlightFields, "...")
		}
		if highlightFields, ok := hit.Highlight["title"]; ok {
			article.Title = strings.Join(highlightFields, "...")
		}
		if highlightFields, ok := hit.Highlight["content_short"]; ok {
			article.ContentShort = strings.Join(highlightFields, "...")
		}

		articles = append(articles, article)
	}

	// 8. 返回结果
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

// 判断知识库类型，选择es索引
func (s *articleService) GetESIndex(ctx context.Context, userId string, kbid uint) string {
	kb, _ := s.kbRepository.GetKBViewById(ctx, kbid)
	index := ""
	if kb.KBType == enums.KBTypePrivate {
		index = enums.Private_knowledge_index + userId
	}
	if kb.KBType == enums.KBTypePublic {
		index = enums.Public_knowledge_index + strconv.Itoa(int(kb.KBID))
	}
	if kb.KBType == enums.KBTypeTeam {
		index = enums.Team_knowledge_index + strconv.Itoa(int(kb.TeamID))
	}
	return index
}
