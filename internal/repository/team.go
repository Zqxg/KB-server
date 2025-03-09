package repository

import (
	"context"
	"go.uber.org/zap"
	"projectName/internal/model"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *model.Team) (uint, error)
	GetTeamByID(ctx context.Context, teamID uint) (*model.Team, error)
	GetTeamList(ctx context.Context, teamName, createdBy string, pageIndex, pageSize int) ([]model.Team, int64, error)
	UpdateTeam(ctx context.Context, team *model.Team) error
	DeleteTeam(ctx context.Context, teamID uint) error

	CreateTeamMember(ctx context.Context, teamMember *model.Member) error
	GetTeamMemberByID(ctx context.Context, teamMemberID uint) (*model.Member, error)
	UpdateTeamMember(ctx context.Context, teamMember *model.Member) error
	DeleteTeamMember(ctx context.Context, teamMemberID uint) error
	UpdateTeamMemberRole(ctx context.Context, teamMemberID uint, role int) error
	GetTeamMemberListByTeamID(ctx context.Context, teamID uint) ([]*model.Member, error)
	GetTeamMemberListByUserID(ctx context.Context, userID string) ([]*model.Member, error)
	GetMemberByTeamIDAndUserID(ctx context.Context, teamID uint, userID string) (*model.Member, error)
}

func NewTeamRepository(
	repository *Repository,
) TeamRepository {
	return &teamRepository{
		Repository: repository,
	}
}

type teamRepository struct {
	*Repository
}

func (r *teamRepository) CreateTeam(ctx context.Context, team *model.Team) (uint, error) {
	if err := r.DB(ctx).Table("sys_team").Create(&team).Error; err != nil {
		r.logger.WithContext(ctx).Error("TeamRepository.CreateTeam error", zap.Error(err))
		return 0, nil
	}
	return team.TeamID, nil
}

func (r *teamRepository) GetTeamByID(ctx context.Context, teamID uint) (*model.Team, error) {
	var team model.Team
	if err := r.DB(ctx).Table("sys_team").Where("team_id = ?", teamID).First(&team).Error; err != nil {
		r.logger.WithContext(ctx).Error("TeamRepository.GetTeamByID error", zap.Error(err))
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) GetTeamList(ctx context.Context, teamName, createdBy string, pageIndex, pageSize int) ([]model.Team, int64, error) {
	var teams []model.Team
	var totalCount int64

	// 构建查询
	db := r.DB(ctx).Table("sys_team").Where("deleted_at IS NULL")

	// 条件判断 - 模糊查询
	if teamName != "" {
		db = db.Where("team_name LIKE ?", "%"+teamName+"%")
	}
	if createdBy != "" {
		db = db.Where("created_by LIKE ?", "%"+createdBy+"%")
	}

	// 分页及数据查询
	if err := db.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&teams).Error; err != nil {
		return nil, 0, err
	}

	return teams, totalCount, nil
}

func (r *teamRepository) UpdateTeam(ctx context.Context, team *model.Team) error {
	if err := r.DB(ctx).Table("sys_team").Where("team_id =?", team.TeamID).Updates(&team).Error; err != nil {
		r.logger.WithContext(ctx).Error("TeamRepository.UpdateTeam error", zap.Error(err))
		return err
	}
	return nil
}
func (r *teamRepository) DeleteTeam(ctx context.Context, teamID uint) error {
	// 软删除团队成员
	if err := r.DB(ctx).Table("sys_member").Where("team_id =?", teamID).Delete(&model.Member{}).Error; err != nil {
		r.logger.WithContext(ctx).Error("TeamRepository.DeleteTeam error", zap.Error(err))
		return err
	}
	// 软删除团队
	if err := r.DB(ctx).Table("sys_team").Where("team_id =?", teamID).Delete(&model.Team{}).Error; err != nil {
		r.logger.WithContext(ctx).Error("TeamRepository.DeleteTeam error", zap.Error(err))
		return err
	}
	return nil
}

func (r *teamRepository) CreateTeamMember(ctx context.Context, teamMember *model.Member) error {
	if err := r.DB(ctx).Table("sys_member").Create(&teamMember).Error; err != nil {
		r.logger.WithContext(ctx).Error("TeamRepository.CreateTeamMember error", zap.Error(err))
		return err
	}
	return nil
}

func (r *teamRepository) GetTeamMemberByID(ctx context.Context, teamMemberID uint) (*model.Member, error) {
	var teamMember model.Member
	if err := r.DB(ctx).Table("sys_member").Where("member_id =?", teamMemberID).First(&teamMember).Error; err != nil {
		r.logger.WithContext(ctx).Error("TeamRepository.GetTeamMemberByID error", zap.Error(err))
		return nil, err
	}
	return &teamMember, nil
}
func (r *teamRepository) UpdateTeamMember(ctx context.Context, teamMember *model.Member) error {
	if err := r.DB(ctx).Table("sys_member").Where("member_id =?", teamMember.MemberID).Updates(&teamMember).Error; err != nil {
		r.logger.WithContext(ctx).Error("TeamRepository.UpdateTeamMember error", zap.Error(err))
		return err
	}
	return nil
}
func (r *teamRepository) DeleteTeamMember(ctx context.Context, teamMemberID uint) error {
	if err := r.DB(ctx).Table("sys_member").Where("member_id =?", teamMemberID).Delete(&model.Member{}).Error; err != nil {
		r.logger.WithContext(ctx).Error("TeamRepository.DeleteTeamMember error", zap.Error(err))
		return err
	}
	return nil
}
func (r *teamRepository) UpdateTeamMemberRole(ctx context.Context, teamMemberID uint, role int) error {
	if err := r.DB(ctx).Table("sys_member").Where("member_id =?", teamMemberID).Update("role", role).Error; err != nil {
		r.logger.WithContext(ctx).Error("TeamRepository.UpdateTeamMemberRole error", zap.Error(err))
		return err
	}
	return nil
}
func (r *teamRepository) GetTeamMemberListByTeamID(ctx context.Context, teamID uint) ([]*model.Member, error) {
	var teamMembers []*model.Member
	if err := r.DB(ctx).Table("sys_member").Where("team_id =?", teamID).Find(&teamMembers).Error; err != nil {
		r.logger.WithContext(ctx).Error("TeamRepository.GetTeamMemberListByTeamID error", zap.Error(err))
		return nil, err
	}
	return teamMembers, nil
}
func (r *teamRepository) GetTeamMemberListByUserID(ctx context.Context, userID string) ([]*model.Member, error) {
	var teamMembers []*model.Member
	if err := r.DB(ctx).Table("sys_member").Where("user_id =?", userID).Find(&teamMembers).Error; err != nil {
		r.logger.WithContext(ctx).Error("TeamRepository.GetTeamMemberListByUserID error", zap.Error(err))
		return nil, err
	}
	return teamMembers, nil
}

func (r *teamRepository) GetMemberByTeamIDAndUserID(ctx context.Context, teamID uint, userID string) (*model.Member, error) {
	var teamMember model.Member
	// `idx_team_user`
	if err := r.DB(ctx).Table("sys_member").Where("team_id =? and user_id =?", teamID, userID).First(&teamMember).Error; err != nil {
		r.logger.WithContext(ctx).Error("TeamRepository.GetMemberByTeamIDAndUserID error", zap.Error(err))
		return nil, err
	}
	return &teamMember, nil
}
