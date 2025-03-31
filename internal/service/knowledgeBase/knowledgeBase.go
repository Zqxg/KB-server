package knowledgeBase

import (
	"errors"
	"github.com/gin-gonic/gin"
	v1 "projectName/api/v1"
	"projectName/internal/enums"
	"projectName/internal/model"
	"projectName/internal/model/vo"
	"projectName/internal/repository"
	"projectName/internal/service"
	"strconv"
)

type KnowledgeBaseService interface {
	// 知识库
	CreateKnowledgeBase(ctx *gin.Context, userId string, req *v1.CreateKBRequest) (uint, error)
	UpdateKBName(ctx *gin.Context, userId string, role int, req *v1.UpdateKBNameReq) error
	DeleteKB(ctx *gin.Context, userId string, role int, kbId uint) error
	GetKBInfo(ctx *gin.Context, kbId uint) (*v1.GetKBInfoResp, error)
	GetKBListByTeamId(ctx *gin.Context, req *v1.GetKBListByTeamIdReq) (*v1.GetKBListByTeamIdResp, error)
	GetKBListByType(ctx *gin.Context, userId string, kbType string) ([]*vo.KbKnowledgeBaseView, error)

	CreateCategory(ctx *gin.Context, userId string, req *v1.CreateCategoryReq) error
	UpdateCategory(ctx *gin.Context, userId string, req *v1.UpdateCategoryReq) error
	DeleteCategory(ctx *gin.Context, userId string, req *v1.DeleteCategoryReq) error
	GetCategoryListByKB(ctx *gin.Context, kbID uint) ([]vo.CategoryView, error)
	//GetCategoryList(ctx *gin.Context, userId string, req *v1.GetCategoryListReq) (*v1.GetCategoryListResp, error)
	//公共知识库
	CreatePublicKB(ctx *gin.Context, userId string, req *v1.CreateKBRequest) (uint, error)
	UpdatePublicKB(ctx *gin.Context, req *v1.UpdateKBNameReq) error
	DeletePublicKB(ctx *gin.Context, kbId uint) error
	//GetPublicKBList(ctx *gin.Context, req *v1.GetKBListRequest) (*v1.GetKBListResponse, error)
	CreatePublicCategory(ctx *gin.Context, req *v1.CreateCategoryReq) error
	UpdatePublicCategory(ctx *gin.Context, req *v1.UpdateCategoryReq) error
	DeletePublicCategory(ctx *gin.Context, req *v1.DeleteCategoryReq) error
}

func NewKnowledgeBaseService(
	service *service.Service,
	kbRepository repository.KBRepository,
	teamRepository repository.TeamRepository,
	articleRepository repository.ArticleRepository,
	userRepository repository.UserRepository,
) KnowledgeBaseService {
	return &knowledgeBaseService{
		Service:           service,
		kbRepository:      kbRepository,
		teamRepository:    teamRepository,
		userRepository:    userRepository,
		articleRepository: articleRepository,
	}
}

type knowledgeBaseService struct {
	*service.Service
	kbRepository      repository.KBRepository
	teamRepository    repository.TeamRepository
	userRepository    repository.UserRepository
	articleRepository repository.ArticleRepository
}

func (s *knowledgeBaseService) CreateKnowledgeBase(ctx *gin.Context, userId string, req *v1.CreateKBRequest) (uint, error) {
	// 校验团队是否存在
	team, err := s.teamRepository.GetTeamByID(ctx, req.TeamID)
	if err != nil {
		return 0, v1.ErrTeamNotExist
	}
	// 获取用户角色
	member, err := s.teamRepository.GetMemberByTeamIDAndUserID(ctx, req.TeamID, userId)
	// 校验用户是否为团队leader或admin
	if err != nil || member.Role != enums.LEADER && member.Role != enums.ADMIN {
		return 0, v1.ErrPermissionDenied
	}
	// 创建知识库
	knowledgeBase := &model.KnowledgeBase{
		TeamID:    team.TeamID,
		KbName:    req.Name,
		CreatedBy: userId,
		UserID:    "",    //团队知识库的userID为空
		IsPublic:  false, //团队知识库
	}
	kbId, err := s.kbRepository.CreateKB(ctx, knowledgeBase)
	if err != nil {
		if errors.Is(err, v1.ErrDuplicateKey) {
			return 0, v1.ErrKnowledgeExist // 自定义错误码
		}
		return 0, v1.ErrCreateKnowledgeFailed
	}
	return kbId, nil
}

func (s *knowledgeBaseService) UpdateKBName(ctx *gin.Context, userId string, role int, req *v1.UpdateKBNameReq) error {
	// 判断知识库是否存在
	kb, err := s.kbRepository.GetKBById(ctx, req.KBID)
	if err != nil {
		return v1.ErrKnowledgeNotExist
	}
	if role == enums.COMMON_USER {
		// 获取用户角色
		member, err := s.teamRepository.GetMemberByTeamIDAndUserID(ctx, kb.TeamID, userId)
		// 校验用户是否为团队leader或admin
		if err != nil || member.Role != enums.LEADER && member.Role != enums.ADMIN {
			return v1.ErrPermissionDenied
		}
	}
	// superAdmin可以直接更新知识库名称
	kb.KbName = req.Name
	err = s.kbRepository.UpdateKB(ctx, kb)
	return v1.ErrUpdateKnowledgeFailed
}

func (s *knowledgeBaseService) DeleteKB(ctx *gin.Context, userId string, role int, kbId uint) error {
	// 判断知识库是否存在
	kb, err := s.kbRepository.GetKBById(ctx, kbId)
	if err != nil {
		return v1.ErrKnowledgeNotExist
	}
	if role == enums.COMMON_USER {
		// 获取用户角色
		member, err := s.teamRepository.GetMemberByTeamIDAndUserID(ctx, kb.TeamID, userId)
		// 校验用户是否为团队leader或admin
		if err != nil || member.Role != enums.LEADER && member.Role != enums.ADMIN {
			return v1.ErrPermissionDenied
		}
	}
	// superAdmin可以直接删除知识库
	err = s.kbRepository.DeleteKB(ctx, kbId)
	if err != nil {
		return v1.ErrDeleteKnowledgeFailed
	}
	return nil
}

func (s *knowledgeBaseService) GetKBInfo(ctx *gin.Context, kbId uint) (*v1.GetKBInfoResp, error) {
	// 判断知识库是否存在
	kb, err := s.kbRepository.GetKBViewById(ctx, kbId)
	if err != nil {
		return nil, v1.ErrKnowledgeNotExist
	}

	// 获取知识库信息
	kbInfo := &v1.GetKBInfoResp{
		KBID:      kb.KBID,
		KbName:    kb.KbName,
		TeamID:    kb.TeamID,
		TeamName:  kb.TeamName,
		IsPublic:  kb.IsPublic,
		UserID:    kb.CreatedBy,
		UserName:  kb.UserName,
		KBType:    kb.KBType,
		CreatedAt: kb.CreatedAt,
		UpdatedAt: kb.UpdatedAt,
	}
	return kbInfo, nil
}

func (s *knowledgeBaseService) GetKBListByTeamId(ctx *gin.Context, req *v1.GetKBListByTeamIdReq) (*v1.GetKBListByTeamIdResp, error) {
	// 初始化分页参数
	pageIndex, pageSize := service.InitPage(req.PageIndex, req.PageSize)

	// 获取知识库列表
	kbList, total, err := s.kbRepository.GetKBListByTeamId(ctx, req.TeamID, pageIndex, pageSize)
	if err != nil {
		return nil, err
	}

	// 组装返回数据
	var kbInfoList []v1.GetKBInfoResp
	for _, kb := range kbList {
		kbInfoList = append(kbInfoList, v1.GetKBInfoResp{
			KBID:      kb.KBID,
			KbName:    kb.KbName,
			TeamID:    kb.TeamID,
			UserID:    kb.UserID,
			TeamName:  kb.TeamName, // 这里假设 `team` 结构体有 `Name` 字段
			UserName:  kb.UserName,
			IsPublic:  kb.IsPublic,
			KBType:    kb.KBType,
			CreatedAt: kb.CreatedAt,
			UpdatedAt: kb.UpdatedAt,
		})
	}

	// 组装响应数据
	resp := &v1.GetKBListByTeamIdResp{
		KBList: kbInfoList,
		PageResponse: v1.PageResponse{
			PageIndex:  pageIndex,
			PageSize:   pageSize,
			TotalCount: total,
		},
	}

	return resp, nil
}

func (s *knowledgeBaseService) CreateCategory(ctx *gin.Context, userId string, req *v1.CreateCategoryReq) error {
	// 校验知识库是否存在
	kbView, err := s.kbRepository.GetKBViewById(ctx, req.KBID)
	if err != nil {
		return v1.ErrKnowledgeNotExist
	}
	// 根据知识库类型新建分类
	if kbView.KBType == enums.KBTypePrivate {
		// 校验知识库userID是否为当前用户
		if kbView.UserID != userId {
			return v1.ErrPermissionDenied
		}
	} else if kbView.KBType == enums.KBTypeTeam {
		// 获取用户角色
		member, err := s.teamRepository.GetMemberByTeamIDAndUserID(ctx, kbView.TeamID, userId)
		// 校验用户是否为团队leader或admin
		if err != nil || member.Role != enums.LEADER && member.Role != enums.ADMIN {
			return v1.ErrPermissionDenied
		}
	} else {
		// 公共知识库因为要校验用户角色，走公共知识库创建分类逻辑
		return v1.ErrCreateCategoryFailed
	}
	// 创建分类
	category := &model.Category{
		KbID:         req.KBID,
		ParentId:     req.ParentId,
		CategoryName: req.CategoryName,
	}
	_, err = s.kbRepository.CreateCategory(ctx, category)
	if err != nil {
		return v1.ErrCreateCategoryFailed
	}
	return nil
}

func (s *knowledgeBaseService) UpdateCategory(ctx *gin.Context, userId string, req *v1.UpdateCategoryReq) error {
	// 判断分类是否存在
	category, err := s.kbRepository.GetCategoryById(ctx, req.CID)
	if err != nil {
		return v1.ErrCategoryNotExist
	}
	// 获取知识库类型
	kbView, err := s.kbRepository.GetKBViewById(ctx, category.KbID)
	if err != nil {
		return v1.ErrKnowledgeNotExist
	}
	// 根据知识库类型更新分类
	if kbView.KBType == enums.KBTypePrivate {
		// 校验知识库userID是否为当前用户
		if kbView.UserID != userId {
			return v1.ErrPermissionDenied
		}
	} else if kbView.KBType == enums.KBTypeTeam {
		// 获取用户角色
		member, err := s.teamRepository.GetMemberByTeamIDAndUserID(ctx, kbView.TeamID, userId)
		// 校验用户是否为团队leader或admin
		if err != nil || member.Role != enums.LEADER && member.Role != enums.ADMIN {
			return v1.ErrPermissionDenied
		}
	} else {
		// 公共知识库因为要校验用户角色，走公共知识库更新分类逻辑
		return v1.ErrUpdateCategoryFailed
	}
	// 更新分类
	category.CategoryName = req.CategoryName
	err = s.kbRepository.UpdateCategory(ctx, category)
	if err != nil {
		return v1.ErrUpdateCategoryFailed
	}
	return nil
}

func (s *knowledgeBaseService) DeleteCategory(ctx *gin.Context, userId string, req *v1.DeleteCategoryReq) error {
	// 判断分类是否存在
	category, err := s.kbRepository.GetCategoryById(ctx, req.CID)
	if err != nil {
		return v1.ErrCategoryNotExist
	}
	// 获取知识库类型
	kbView, err := s.kbRepository.GetKBViewById(ctx, category.KbID)
	if err != nil {
		return v1.ErrKnowledgeNotExist
	}
	// 根据知识库类型删除分类
	if kbView.KBType == enums.KBTypePrivate {
		// 校验知识库userID是否为当前用户
		if kbView.UserID != userId {
			return v1.ErrPermissionDenied
		}
	} else if kbView.KBType == enums.KBTypeTeam {
		// 获取用户角色
		member, err := s.teamRepository.GetMemberByTeamIDAndUserID(ctx, kbView.TeamID, userId)
		// 校验用户是否为团队leader或admin
		if err != nil || member.Role != enums.LEADER && member.Role != enums.ADMIN {
			return v1.ErrPermissionDenied
		}
	} else {
		// 公共知识库因为要校验用户角色，走公共知识库删除分类逻辑
		return v1.ErrDeleteCategoryFailed
	}
	// 删除分类
	err = s.kbRepository.DeleteCategory(ctx, req.CID)
	if err != nil {
		return v1.ErrDeleteCategoryFailed
	}
	return nil
}

func (s *knowledgeBaseService) GetCategoryListByKB(ctx *gin.Context, kbID uint) ([]vo.CategoryView, error) {
	// 获取分类列表
	categoryList, err := s.kbRepository.GetCategoryTreeByKB(ctx, kbID)
	if err != nil {
		return nil, v1.ErrGetCategoryListFailed
	}
	return categoryList, nil
}

func (s *knowledgeBaseService) GetKBListByType(ctx *gin.Context, userId string, kbType string) ([]*vo.KbKnowledgeBaseView, error) {

	// 根据类型获取知识库列表
	if kbType == "public" {
		kbList, err := s.kbRepository.GetKBListByTypeAndUserId(ctx, "", "publicKB")
		if err != nil {
			return nil, v1.ErrKnowledgeNotExist
		}
		return kbList, nil
	}
	if kbType == "private" {
		kbList, err := s.kbRepository.GetKBListByTypeAndUserId(ctx, userId, "privateKB")
		if err != nil {
			return nil, v1.ErrKnowledgeNotExist
		}
		return kbList, nil
	}
	return nil, v1.ErrKnowledgeNotExist
}

func (s *knowledgeBaseService) CreatePublicKB(ctx *gin.Context, userId string, req *v1.CreateKBRequest) (uint, error) {
	// 创建知识库
	knowledgeBase := &model.KnowledgeBase{
		KbName:    req.Name,
		CreatedBy: userId,
		IsPublic:  true, //团队知识库
	}
	kbId, err := s.kbRepository.CreateKB(ctx, knowledgeBase)
	if err != nil {
		return 0, err
	}
	// 新增ES索引 公共知识库
	index := enums.Public_knowledge_index + strconv.Itoa(int(kbId))
	err = s.articleRepository.CreateEsIndex(ctx, index)
	if err != nil {
		return 0, v1.ErrCreateEsIndexFailed
	}
	return kbId, nil
}

func (s *knowledgeBaseService) UpdatePublicKB(ctx *gin.Context, req *v1.UpdateKBNameReq) error {
	// 判断知识库是否存在
	kb, err := s.kbRepository.GetKBById(ctx, req.KBID)
	if err != nil {
		return v1.ErrKnowledgeNotExist
	}
	// 更新知识库名称
	kb.KbName = req.Name
	err = s.kbRepository.UpdateKB(ctx, kb)
	if err != nil {
		return v1.ErrUpdateKnowledgeFailed
	}
	return nil
}

func (s *knowledgeBaseService) DeletePublicKB(ctx *gin.Context, kbId uint) error {
	// 判断知识库是否存在
	_, err := s.kbRepository.GetKBById(ctx, kbId)
	if err != nil {
		return v1.ErrKnowledgeNotExist
	}
	// 删除ES索引
	index := enums.Public_knowledge_index + strconv.Itoa(int(kbId))
	err = s.articleRepository.DeleteEsIndex(ctx, index)
	if err != nil {
		return v1.ErrDeleteEsIndexFailed
	}
	// 删除知识库
	err = s.kbRepository.DeleteKB(ctx, kbId)
	if err != nil {
		return v1.ErrDeleteKnowledgeFailed
	}
	return nil
}

func (s *knowledgeBaseService) CreatePublicCategory(ctx *gin.Context, req *v1.CreateCategoryReq) error {
	// 校验知识库是否存在
	kb, err := s.kbRepository.GetKBById(ctx, req.KBID)
	if err != nil {
		return v1.ErrKnowledgeNotExist
	}
	// 校验知识库是否为公共知识库
	if !kb.IsPublic {
		return v1.ErrNotPublicKnowledge
	}
	// 创建分类
	category := &model.Category{
		KbID:         req.KBID,
		ParentId:     req.ParentId,
		CategoryName: req.CategoryName,
	}
	_, err = s.kbRepository.CreateCategory(ctx, category)
	if err != nil {
		return v1.ErrCreateCategoryFailed
	}
	return nil
}

func (s *knowledgeBaseService) UpdatePublicCategory(ctx *gin.Context, req *v1.UpdateCategoryReq) error {
	// 判断分类是否存在
	category, err := s.kbRepository.GetCategoryById(ctx, req.CID)
	if err != nil {
		return v1.ErrCategoryNotExist
	}
	// 校验知识库是否为公共知识库
	kb, err := s.kbRepository.GetKBById(ctx, category.KbID)
	if err != nil {
		return v1.ErrKnowledgeNotExist
	}
	if !kb.IsPublic {
		return v1.ErrNotPublicKnowledge
	}
	// 更新分类
	category.CategoryName = req.CategoryName
	err = s.kbRepository.UpdateCategory(ctx, category)
	if err != nil {
		return v1.ErrUpdateCategoryFailed
	}
	return nil
}
func (s *knowledgeBaseService) DeletePublicCategory(ctx *gin.Context, req *v1.DeleteCategoryReq) error {
	// 判断分类是否存在
	category, err := s.kbRepository.GetCategoryById(ctx, req.CID)
	if err != nil {
		return v1.ErrCategoryNotExist
	}
	// 校验知识库是否为公共知识库
	kb, err := s.kbRepository.GetKBById(ctx, category.KbID)
	if err != nil {
		return v1.ErrKnowledgeNotExist
	}
	if !kb.IsPublic {
		return v1.ErrNotPublicKnowledge
	}
	// 删除分类
	err = s.kbRepository.DeleteCategory(ctx, req.CID)
	if err != nil {
		return v1.ErrDeleteCategoryFailed
	}
	return nil
}
