package team

import (
	"github.com/gin-gonic/gin"
	v1 "projectName/api/v1"
	"projectName/internal/enums"
	"projectName/internal/model"
	"projectName/internal/repository"
	"projectName/internal/service"
	"strconv"
)

type TeamService interface {
	CreateTeam(ctx *gin.Context, userId string, req *v1.CreateTeamRequest) (uint, error)
	UpdateTeam(ctx *gin.Context, userId string, req *v1.UpdateTeamRequest) error
	DeleteTeam(ctx *gin.Context, userId string, teamId uint) error
	GetTeamList(ctx *gin.Context, req *v1.GetTeamListReq) (*v1.GetTeamListResp, error)
	GetTeamInfo(ctx *gin.Context, teamId int) (*v1.GetTeamInfoResp, error)
	GetUserTeamList(ctx *gin.Context, userId string, pageIndex, pageSize int) (*v1.GetUserTeamListResp, error)
	GetTeamMemberList(ctx *gin.Context, teamId int) ([]*v1.MemberData, error)
	AddTeamMember(ctx *gin.Context, userID string, req *v1.AddTeamMemberReq) error
	DeleteTeamMember(ctx *gin.Context, userID string, req *v1.DeleteTeamMemberReq) error
	UpdateTeamMemberRole(ctx *gin.Context, userID string, req *v1.UpdateTeamMemberRoleReq) error
	QuitTeam(ctx *gin.Context, userID string, teamId uint) error
}

func NewTeamService(
	service *service.Service,
	userRepository repository.UserRepository,
	articleRepository repository.ArticleRepository,
	teamRepository repository.TeamRepository,
) TeamService {
	return &teamService{
		Service:           service,
		userRepository:    userRepository,
		teamRepository:    teamRepository,
		articleRepository: articleRepository,
	}
}

type teamService struct {
	*service.Service
	userRepository    repository.UserRepository
	teamRepository    repository.TeamRepository
	articleRepository repository.ArticleRepository
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
			return 0, v1.ErrCreateTeamFailed
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
		return v1.ErrUpdateTeamFailed
	}
	return nil
}

func (s *teamService) DeleteTeam(ctx *gin.Context, userId string, teamId uint) error {
	// 判断用户是否为团队负责人或者管理员
	if !s.isTeamLeaderOrAdmin(ctx, teamId, userId) {
		return v1.ErrPermissionDenied
	}
	// 删除es索引团队
	index := enums.Team_knowledge_index + strconv.Itoa(int(teamId))
	err := s.articleRepository.DeleteEsIndex(ctx, index)
	if err != nil {
		return v1.ErrDeleteEsIndexFailed
	}
	// 删除团队
	if err := s.teamRepository.DeleteTeam(ctx, teamId); err != nil {
		return v1.ErrDeleteTeamFailed
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

func (s *teamService) GetTeamList(ctx *gin.Context, req *v1.GetTeamListReq) (*v1.GetTeamListResp, error) {
	// 初始化分页信息
	pageIndex, pageSize := service.InitPage(req.PageIndex, req.PageSize)

	// 查询团队列表
	teams, totalCount, err := s.teamRepository.GetTeamList(ctx, req.TeamName, req.CreatedBy, pageIndex, pageSize)
	if err != nil {
		return nil, err
	}
	var teamList []v1.TeamData
	// 转换团队信息
	for _, team := range teams {
		teamList = append(teamList, v1.TeamData{
			TeamID:      team.TeamID,
			TeamName:    team.TeamName,
			Description: team.Description,
			CreatedBy:   team.CreatedBy,
			CreatedAt:   team.CreatedAt,
			UpdatedAt:   team.UpdatedAt,
		})
	}
	// 构建响应
	return &v1.GetTeamListResp{
		TeamList: teamList,
		PageResponse: v1.PageResponse{
			TotalCount: totalCount, // 增加总数，方便前端分页
			PageIndex:  pageIndex,
			PageSize:   pageSize,
		},
	}, nil
}

func (s *teamService) GetTeamInfo(ctx *gin.Context, teamId int) (*v1.GetTeamInfoResp, error) {
	team, err := s.teamRepository.GetTeamByID(ctx, uint(teamId))
	if err != nil {
		return nil, v1.ErrGetTeamInfoFailed
	}
	members, err := s.teamRepository.GetTeamMemberListByTeamID(ctx, uint(teamId))
	if err != nil {
		return nil, v1.ErrGetTeamMemberListFailed
	}

	// 构建响应
	teamData := v1.TeamData{
		TeamID:      team.TeamID,
		TeamName:    team.TeamName,
		Description: team.Description,
		CreatedBy:   team.CreatedBy,
		CreatedAt:   team.CreatedAt,
		UpdatedAt:   team.UpdatedAt,
	}
	var membersData []*v1.MemberData
	for _, member := range members {
		user, err := s.userRepository.GetByUserId(ctx, member.UserID)
		if err != nil {
			return nil, v1.ErrMemberNotExist
		}
		membersData = append(membersData, &v1.MemberData{
			MemberID:   member.MemberID,
			UserID:     member.UserID,
			RoleType:   member.Role,
			JoinTime:   member.CreatedAt,
			UpdateTime: member.UpdatedAt,
			NickName:   user.Nickname,
		})
	}
	return &v1.GetTeamInfoResp{
		Team:   teamData,
		Member: membersData,
	}, nil

}

func (s *teamService) GetUserTeamList(ctx *gin.Context, userId string, pageIndex, pageSize int) (*v1.GetUserTeamListResp, error) {
	index, size := service.InitPage(pageIndex, pageSize)
	teams, total, err := s.teamRepository.GetTeamList(ctx, "", userId, index, size)
	if err != nil {
		return nil, v1.ErrGetTeamInfoFailed
	}
	var teamList []v1.TeamData
	// 转换团队信息
	for _, team := range teams {
		teamList = append(teamList, v1.TeamData{
			TeamID:      team.TeamID,
			TeamName:    team.TeamName,
			Description: team.Description,
			CreatedBy:   team.CreatedBy,
			CreatedAt:   team.CreatedAt,
			UpdatedAt:   team.UpdatedAt,
		})
	}
	return &v1.GetUserTeamListResp{
		TeamList: teamList,
		PageResponse: v1.PageResponse{
			TotalCount: total,
			PageIndex:  index,
			PageSize:   size,
		},
	}, nil
}

func (s *teamService) GetTeamMemberList(ctx *gin.Context, teamId int) ([]*v1.MemberData, error) {
	// 判断团队是否存在
	team, err := s.teamRepository.GetTeamByID(ctx, uint(teamId))
	if err != nil || team == nil {
		return nil, v1.ErrTeamNotExist
	}

	memberList, err := s.teamRepository.GetTeamMemberListByTeamID(ctx, uint(teamId))
	if err != nil {
		return nil, v1.ErrGetTeamMemberListFailed
	}
	var memberDataList []*v1.MemberData
	for _, member := range memberList {
		user, err := s.userRepository.GetByUserId(ctx, member.UserID)
		if err != nil {
			return nil, v1.ErrMemberNotExist
		}
		memberDataList = append(memberDataList, &v1.MemberData{
			MemberID:   member.MemberID,
			UserID:     member.UserID,
			RoleType:   member.Role,
			JoinTime:   member.CreatedAt,
			UpdateTime: member.UpdatedAt,
			NickName:   user.Nickname,
		})
	}
	return memberDataList, nil
}

func (s *teamService) AddTeamMember(ctx *gin.Context, userID string, req *v1.AddTeamMemberReq) error {
	// 判断用户是否为团队负责人或者管理员
	if !s.isTeamLeaderOrAdmin(ctx, req.TeamID, userID) {
		return v1.ErrNoTeamAdminPermission
	}
	// 判断用户是否存在
	_, err := s.userRepository.GetByUserId(ctx, req.MemberID)
	if err != nil {
		return v1.ErrUserNotExist
	}
	// 判断用户是否已经在团队中
	_, err = s.teamRepository.GetMemberByTeamIDAndUserID(ctx, req.TeamID, req.MemberID)
	if err != nil {
		return v1.ErrMemberExist
	}
	// 构建成员表
	teamMember := &model.Member{
		TeamID: req.TeamID,
		UserID: req.MemberID,
		Role:   req.RoleType,
	}
	err = s.teamRepository.CreateTeamMember(ctx, teamMember)
	if err != nil {
		return err
	}
	return nil
}

func (s *teamService) DeleteTeamMember(ctx *gin.Context, userID string, req *v1.DeleteTeamMemberReq) error {
	// 判断用户是否为团队负责人或者管理员
	if !s.isTeamLeaderOrAdmin(ctx, req.TeamID, userID) {
		return v1.ErrNoTeamAdminPermission
	}
	// 判断用户是否在团队中
	member, err := s.teamRepository.GetMemberByTeamIDAndUserID(ctx, req.TeamID, req.MemberID)
	if err != nil {
		return v1.ErrMemberNotExist
	}
	// 删除成员
	err = s.teamRepository.DeleteTeamMember(ctx, member.MemberID)
	if err != nil {
		return err
	}
	return nil
}

func (s *teamService) UpdateTeamMemberRole(ctx *gin.Context, userID string, req *v1.UpdateTeamMemberRoleReq) error {
	// 判断用户是否为团队负责人或者管理员
	if !s.isTeamLeaderOrAdmin(ctx, req.TeamID, userID) {
		return v1.ErrNoTeamAdminPermission
	}
	// 判断用户是否在团队中
	member, err := s.teamRepository.GetMemberByTeamIDAndUserID(ctx, req.TeamID, req.MemberID)
	if err != nil {
		return v1.ErrMemberNotExist
	}
	// 更新成员角色
	member.Role = req.RoleType
	err = s.teamRepository.UpdateTeamMember(ctx, member)
	if err != nil {
		return v1.ErrUpdateMemberFailed
	}
	return nil
}

func (s *teamService) QuitTeam(ctx *gin.Context, userID string, teamId uint) error {
	// 判断用户是否在团队中
	member, err := s.teamRepository.GetMemberByTeamIDAndUserID(ctx, teamId, userID)
	if err != nil {
		return v1.ErrMemberNotExist
	}
	// 删除成员
	err = s.teamRepository.DeleteTeamMember(ctx, member.MemberID)
	if err != nil {
		return v1.ErrDeleteMemberFailed
	}
	return nil
}
