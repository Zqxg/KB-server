package team

import (
	"github.com/gin-gonic/gin"
	v1 "projectName/api/v1"
	"projectName/internal/enums"
	"projectName/internal/model"
	"projectName/internal/repository"
	"projectName/internal/service"
)

type TeamService interface {
	CreateTeam(ctx *gin.Context, userId string, req *v1.CreateTeamRequest) (uint, error)
	UpdateTeam(ctx *gin.Context, userId string, req *v1.UpdateTeamRequest) error
	DeleteTeam(ctx *gin.Context, userId string, teamId uint) error
	GetTeamList(ctx *gin.Context, req *v1.GetTeamListReq) (v1.GetTeamListResp, error)
}

func NewTeamService(
	service *service.Service,
	teamRepository repository.TeamRepository,
) TeamService {
	return &teamService{
		Service:        service,
		teamRepository: teamRepository,
	}
}

type teamService struct {
	*service.Service
	teamRepository repository.TeamRepository
}

func (s *teamService) CreateTeam(ctx *gin.Context, userId string, req *v1.CreateTeamRequest) (uint, error) {
	// 构建 Team 结构体
	team := &model.Team{
		TeamName:    req.TeamName,
		Description: req.Description,
		CreatedBy:   userId,
	}
	teamId, err := s.teamRepository.CreateTeam(ctx, team)
	if err != nil {
		return 0, err
	}
	if teamId != 0 {
		// 构建成员表
		teamMember := &model.Member{
			TeamID: teamId,
			UserID: userId,
			Role:   enums.LEADER, // 负责人
		}
		err = s.teamRepository.CreateTeamMember(ctx, teamMember)
		if err != nil {
			return 0, err
		}
	}
	return teamId, nil
}

func (s *teamService) UpdateTeam(ctx *gin.Context, userId string, req *v1.UpdateTeamRequest) error {
	// 判断用户是否为团队负责人或者管理员
	if !s.isTeamLeaderOrAdmin(ctx, req.TeamID, userId) {
		return v1.ErrPermissionDenied
	}
	// 更新团队信息
	team := &model.Team{
		TeamID:      req.TeamID,
		TeamName:    req.TeamName,
		Description: req.Description,
	}
	if err := s.teamRepository.UpdateTeam(ctx, team); err != nil {
		return err
	}
	return nil
}

func (s *teamService) DeleteTeam(ctx *gin.Context, userId string, teamId uint) error {
	// 判断用户是否为团队负责人或者管理员
	if !s.isTeamLeaderOrAdmin(ctx, teamId, userId) {
		return v1.ErrPermissionDenied
	}
	// 删除团队
	if err := s.teamRepository.DeleteTeam(ctx, teamId); err != nil {
		return err
	}
	return nil
}

// 判断用户是否为团队负责人或者管理员
func (s *teamService) isTeamLeaderOrAdmin(ctx *gin.Context, teamId uint, userId string) bool {
	// 获取用户角色
	member, err := s.teamRepository.GetMemberByTeamIDAndUserID(ctx, teamId, userId)
	if err != nil {
		return false
	}
	if member == nil {
		return false
	}
	return member.Role == enums.LEADER || member.Role == enums.ADMIN
}

func (s *teamService) GetTeamList(ctx *gin.Context, req *v1.GetTeamListReq) (v1.GetTeamListResp, error) {
	// 初始化分页信息
	pageIndex, pageSize := service.InitPage(req.PageIndex, req.PageSize)

	// 查询团队列表
	teams, totalCount, err := s.teamRepository.GetTeamList(ctx, req.TeamName, req.CreatedBy, pageIndex, pageSize)
	if err != nil {
		return v1.GetTeamListResp{}, err
	}

	// 构建响应
	return v1.GetTeamListResp{
		TeamList: teams,
		PageResponse: v1.PageResponse{
			TotalCount: totalCount, // 增加总数，方便前端分页
			PageIndex:  pageIndex,
			PageSize:   pageSize,
		},
	}, nil
}
