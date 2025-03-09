package repository

import (
	"context"
	"go.uber.org/zap"
	"projectName/internal/model"
)

type KBRepository interface {
	CreateKB(ctx context.Context, knowledge *model.KnowledgeBase) (uint, error)
	UpdateKB(ctx context.Context, knowledge *model.KnowledgeBase) error
	DeleteKB(ctx context.Context, id int) error
	DeleteKBByUserId(ctx context.Context, userId string) error
	GetKnowledgeBaseById(ctx context.Context, id int) (*model.KnowledgeBase, error)
	GetKBListByTeamId(ctx context.Context, teamId int) ([]*model.KnowledgeBase, int64, error)
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

func (r *kbRepository) DeleteKB(ctx context.Context, id int) error {
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

func (r *kbRepository) GetKnowledgeBaseById(ctx context.Context, id int) (*model.KnowledgeBase, error) {
	var knowledge model.KnowledgeBase
	if err := r.DB(ctx).Table("kb_knowledge").Where("kb_id =?", id).First(&knowledge).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.GetKB error", zap.Error(err))
		return nil, err
	}
	return &knowledge, nil
}

// GetKBListByTeamId todo 修改成视图
func (r *kbRepository) GetKBListByTeamId(ctx context.Context, teamId int) ([]*model.KnowledgeBase, int64, error) {
	var knowledge []*model.KnowledgeBase
	var count int64
	if err := r.DB(ctx).Table("kb_knowledge").Where("team_id =?", teamId).Count(&count).Find(&knowledge).Error; err != nil {
		r.logger.WithContext(ctx).Error("KBRepository.GetKBList error", zap.Error(err))
		return nil, 0, err
	}
	return knowledge, count, nil
}
