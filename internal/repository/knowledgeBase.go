package repository

import (
	"context"
	"go.uber.org/zap"
	"projectName/internal/model"
	"projectName/internal/model/vo"
)

type KBRepository interface {
	CreateKB(ctx context.Context, knowledge *model.KnowledgeBase) (uint, error)
	UpdateKB(ctx context.Context, knowledge *model.KnowledgeBase) error
	DeleteKB(ctx context.Context, id uint) error
	DeleteKBByUserId(ctx context.Context, userId string) error
	GetKBById(ctx context.Context, id uint) (*model.KnowledgeBase, error)
	GetKBViewById(ctx context.Context, id uint) (*vo.KbKnowledgeBaseView, error)
	GetKBListByTeamId(ctx context.Context, teamId uint, pageIndex, pageSize int) ([]*vo.KbKnowledgeBaseView, int64, error)
	//GetKBList(ctx context.Context, req *model.GetKBListReq) ([]*model.KnowledgeBase, int64, error)
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
	if err := r.DB(ctx).Table("kb_knowledge").Create(&knowledge).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.CreateKB error", zap.Error(err))
		return 0, err
	}
	return knowledge.KbID, nil
}

func (r *kbRepository) UpdateKB(ctx context.Context, knowledge *model.KnowledgeBase) error {
	if err := r.DB(ctx).Table("kb_knowledge").Where("kb_id = ?", knowledge.KbID).Updates(&knowledge).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.UpdateKB error", zap.Error(err))
		return err
	}
	return nil
}

func (r *kbRepository) DeleteKB(ctx context.Context, id uint) error {
	if err := r.DB(ctx).Table("kb_knowledge").Where("kb_id =?", id).Delete(&model.KnowledgeBase{}).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.DeleteKB error", zap.Error(err))
		return err
	}
	return nil
}

func (r *kbRepository) DeleteKBByUserId(ctx context.Context, userId string) error {
	if err := r.DB(ctx).Table("kb_knowledge").Where("user_id =?", userId).Delete(&model.KnowledgeBase{}).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.DeleteKB error", zap.Error(err))
		return err
	}
	return nil
}
func (r *kbRepository) GetKBById(ctx context.Context, id uint) (*model.KnowledgeBase, error) {
	var knowledge model.KnowledgeBase
	if err := r.DB(ctx).Table("kb_knowledge").Where("kb_id =?", id).First(&knowledge).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.GetKB error", zap.Error(err))
		return nil, err
	}
	return &knowledge, nil
}

func (r *kbRepository) GetKBViewById(ctx context.Context, id uint) (*vo.KbKnowledgeBaseView, error) {
	var knowledge vo.KbKnowledgeBaseView
	if err := r.DB(ctx).Table("kb_knowledgeBase_view").Where("kb_id =?", id).First(&knowledge).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.GetKB error", zap.Error(err))
		return nil, err
	}
	return &knowledge, nil
}

// GetKBListByTeamId 获取指定 teamId 的知识库列表（分页）
func (r *kbRepository) GetKBListByTeamId(ctx context.Context, teamId uint, pageIndex, pageSize int) ([]*vo.KbKnowledgeBaseView, int64, error) {
	var kbList []*vo.KbKnowledgeBaseView
	var total int64

	db := r.DB(ctx).Table("kb_knowledgeBase_view").Where("team_id = ?", teamId)

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
