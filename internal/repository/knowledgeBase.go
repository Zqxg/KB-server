package repository

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	v1 "projectName/api/v1"
	"projectName/internal/model"
	"projectName/internal/model/vo"
)

type KBRepository interface {
	CreateKB(ctx context.Context, knowledge *model.KnowledgeBase) (uint, error)
	UpdateKB(ctx context.Context, knowledge *model.KnowledgeBase) error
	DeleteKB(ctx context.Context, id uint) error
	DeleteKBByTeamID(ctx context.Context, teamID uint) error
	DeleteKBByUserId(ctx context.Context, userId string) error
	GetKBById(ctx context.Context, id uint) (*model.KnowledgeBase, error)
	GetKBViewById(ctx context.Context, id uint) (*vo.KbKnowledgeBaseView, error)
	GetKBTypeById(ctx context.Context, id uint) (string, error)
	GetKBListByTeamId(ctx context.Context, teamId uint, pageIndex, pageSize int) ([]*vo.KbKnowledgeBaseView, int64, error)
	GetKBListByTypeAndUserId(ctx context.Context, userId, kbType string) ([]*vo.KbKnowledgeBaseView, error)
	//GetKBList(ctx context.Context, req *model.GetKBListReq) ([]*model.KnowledgeBase, int64, error)
	CreateCategory(ctx context.Context, category *model.Category) (uint, error)
	UpdateCategory(ctx context.Context, category *model.Category) error
	DeleteCategory(ctx context.Context, id uint) error
	GetCategoryList(ctx context.Context, kbId uint) ([]*model.Category, error)
	GetCategoryById(ctx context.Context, id uint) (*model.Category, error)
	GetCategoryTreeByKB(ctx context.Context, kbId uint) ([]vo.CategoryView, error)
}

func NewKBRepository(
	repository *Repository,
) KBRepository {
	return &kbRepository{
		Repository: repository,
	}
}

type kbRepository struct {
	*Repository
}

func (r *kbRepository) CreateKB(ctx context.Context, knowledge *model.KnowledgeBase) (uint, error) {
	var existing model.KnowledgeBase
	err := r.DB(ctx).Table("kb_knowledgeBase").
		Where("kb_name = ? AND user_id = ? AND deleted_at IS NULL", knowledge.KbName, knowledge.UserID).
		First(&existing).Error

	if err == nil {
		return 0, v1.ErrDuplicateKey // 说明已存在未删除的记录
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err // 其他数据库错误
	}

	// 继续插入新记录
	if err := r.DB(ctx).Table("kb_knowledgeBase").Create(&knowledge).Error; err != nil {
		return 0, err
	}
	return knowledge.KbID, nil
}

func (r *kbRepository) UpdateKB(ctx context.Context, knowledge *model.KnowledgeBase) error {
	if err := r.DB(ctx).Table("kb_knowledgeBase").Where("kb_id = ?", knowledge.KbID).Updates(&knowledge).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.UpdateKB error", zap.Error(err))
		return err
	}
	return nil
}

// 公共删除方法
func (r *kbRepository) deleteKBs(ctx context.Context, kbIDs []uint) error {
	if len(kbIDs) == 0 {
		return nil
	}

	// 事务保证数据一致性
	return r.DB(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 删除知识库
		if err := tx.Table("kb_knowledgeBase").Where("kb_id IN ?", kbIDs).Delete(&model.KnowledgeBase{}).Error; err != nil {
			return err
		}
		// 2. 删除分类
		if err := tx.Table("kb_category").Where("kb_id IN ?", kbIDs).Delete(&model.Category{}).Error; err != nil {
			return err
		}
		// 3. 删除文章
		if err := tx.Table("kb_article").Where("kb_id IN ?", kbIDs).Delete(&model.Article{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteKB 根据 kb_id 删除单个知识库
func (r *kbRepository) DeleteKB(ctx context.Context, id uint) error {
	return r.deleteKBs(ctx, []uint{id})
}

// DeleteKBByTeamID 根据 team_id 删除该团队下的所有知识库
func (r *kbRepository) DeleteKBByTeamID(ctx context.Context, teamID uint) error {
	var kbIDs []uint
	if err := r.DB(ctx).Table("kb_knowledgeBase").Where("team_id = ?", teamID).Pluck("kb_id", &kbIDs).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.GetKBIDs error", zap.Error(err))
		return err
	}
	return r.deleteKBs(ctx, kbIDs)
}

// DeleteKBByUserId 根据 user_id 删除该用户创建的所有知识库
func (r *kbRepository) DeleteKBByUserId(ctx context.Context, userId string) error {
	var kbIDs []uint
	if err := r.DB(ctx).Table("kb_knowledgeBase").Where("user_id = ?", userId).Pluck("kb_id", &kbIDs).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.GetKBIDs error", zap.Error(err))
		return err
	}
	return r.deleteKBs(ctx, kbIDs)
}

func (r *kbRepository) GetKBById(ctx context.Context, id uint) (*model.KnowledgeBase, error) {
	var knowledge model.KnowledgeBase
	if err := r.DB(ctx).Table("kb_knowledgeBase").Where("kb_id =?", id).First(&knowledge).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.GetKB error", zap.Error(err))
		return nil, err
	}
	return &knowledge, nil
}

func (r *kbRepository) GetKBViewById(ctx context.Context, id uint) (*vo.KbKnowledgeBaseView, error) {
	var knowledge vo.KbKnowledgeBaseView
	if err := r.DB(ctx).Table("kb_knowledge_base_view").Where("kb_id =?", id).First(&knowledge).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.GetKB error", zap.Error(err))
		return nil, err
	}
	return &knowledge, nil
}
func (r *kbRepository) GetKBTypeById(ctx context.Context, id uint) (string, error) {
	var knowledgeView *vo.KbKnowledgeBaseView
	if err := r.DB(ctx).Table("kb_knowledge_base_view").Where("kb_id =?", id).First(&knowledgeView).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.GetKBTypeById error", zap.Error(err))
		return "", err
	}
	return knowledgeView.KBType, nil
}

// GetKBListByTeamId 获取指定 teamId 的知识库列表（分页）
func (r *kbRepository) GetKBListByTeamId(ctx context.Context, teamId uint, pageIndex, pageSize int) ([]*vo.KbKnowledgeBaseView, int64, error) {
	var kbList []*vo.KbKnowledgeBaseView
	var total int64

	db := r.DB(ctx).Table("kb_knowledge_base_view").Where("team_id = ?", teamId)

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.GetKBListByTeamId count error", zap.Error(err))
		return nil, 0, err
	}

	// 分页查询
	offset := (pageIndex - 1) * pageSize
	if err := db.Limit(pageSize).Offset(offset).Find(&kbList).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.GetKBListByTeamId query error", zap.Error(err))
		return nil, 0, err
	}

	return kbList, total, nil
}

// CreateCategory 创建分类
func (r *kbRepository) CreateCategory(ctx context.Context, category *model.Category) (uint, error) {
	if err := r.DB(ctx).Table("kb_category").Create(&category).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.CreateCategory error", zap.Error(err))
		return 0, err
	}
	return category.CategoryId, nil
}

// UpdateCategory 更新分类
func (r *kbRepository) UpdateCategory(ctx context.Context, category *model.Category) error {
	if err := r.DB(ctx).Table("kb_category").Where("category_id =?", category.CategoryId).Updates(&category).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.UpdateCategory error", zap.Error(err))
		return err
	}
	return nil
}

// DeleteCategory 删除分类
func (r *kbRepository) DeleteCategory(ctx context.Context, id uint) error {
	if err := r.DB(ctx).Table("kb_category").Where("category_id =?", id).Delete(&model.Category{}).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.DeleteCategory error", zap.Error(err))
		return err
	}
	return nil
}

// GetCategoryList 获取分类列表
func (r *kbRepository) GetCategoryList(ctx context.Context, kbId uint) ([]*model.Category, error) {
	var categories []*model.Category
	if err := r.DB(ctx).Table("kb_category").Where("kb_id =?", kbId).Find(&categories).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.GetCategoryList error", zap.Error(err))
		return nil, err
	}
	return categories, nil
}

// GetCategoryById 获取分类信息
func (r *kbRepository) GetCategoryById(ctx context.Context, id uint) (*model.Category, error) {
	var category model.Category
	if err := r.DB(ctx).Table("kb_category").Where("category_id =?", id).First(&category).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.GetCategoryById error", zap.Error(err))
		return nil, err
	}
	return &category, nil
}

// GetCategoryTreeByKB 获取指定知识库的分类树结构
func (r *kbRepository) GetCategoryTreeByKB(ctx context.Context, kbId uint) ([]vo.CategoryView, error) {
	var categories []vo.CategoryView
	if err := r.DB(ctx).Table("kb_category_view").
		Where("kb_id = ?", kbId).
		Order("level ASC, category_id ASC"). // 先按层级、再按ID排序，方便递归
		Find(&categories).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.GetCategoryTreeByKB error", zap.Error(err))
		return nil, err
	}

	tree := buildCategoryTree(categories, 0) // 从 parent_id = 0 开始构建树
	return tree, nil
}

// buildCategoryTree 递归生成分类树
func buildCategoryTree(data []vo.CategoryView, parentId uint) []vo.CategoryView {
	var result []vo.CategoryView
	for _, item := range data {
		if item.ParentId == parentId {
			item.Children = buildCategoryTree(data, item.CId)
			result = append(result, item)
		}
	}
	return result
}

func (r *kbRepository) GetKBListByTypeAndUserId(ctx context.Context, userId, kbType string) ([]*vo.KbKnowledgeBaseView, error) {
	var kbList []*vo.KbKnowledgeBaseView
	query := r.DB(ctx).Table("kb_knowledge_base_view").Where("kb_type = ?", kbType)

	// 如果 userId 是空字符串，则查询 NULL
	query = query.Where("user_id = ?", userId)

	if err := query.Find(&kbList).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.GetKBListByTypeAndUserId error", zap.Error(err))
		return nil, err
	}
	return kbList, nil
}
