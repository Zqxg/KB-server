package team

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	ApplyJoinTeam(ctx *gin.Context, userID string, req *v1.ApplyJoinTeamReq) (uint, error)
}

func NewTeamService(
	service *service.Service,
	userRepository repository.UserRepository,
	articleRepository repository.ArticleRepository,
	kbRepository repository.KBRepository,
	teamRepository repository.TeamRepository,
) TeamService {
	return &teamService{
		Service:           service,
		userRepository:    userRepository,
		teamRepository:    teamRepository,
		articleRepository: articleRepository,
		kbRepository:      kbRepository,
	}
}

type teamService struct {
	*service.Service
	userRepository    repository.UserRepository
	teamRepository    repository.TeamRepository
	articleRepository repository.ArticleRepository
	kbRepository      repository.KBRepository
}

func (s *teamService) CreateTeam(ctx *gin.Context, userId string, req *v1.CreateTeamRequest) (uint, error) {
	// 判断团队名称长度
	if len(req.TeamName) < 2 {
		return 0, v1.ErrTeamNameTooShort
	}
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
		// 创建es索引
		index := enums.Team_knowledge_index + strconv.Itoa(int(teamId))
		esMapper, err := service.GenerateEsMapping(model.EsArticle{})
		if err != nil {
			return 0, v1.ErrCreateEsMapperFailed
		}
		err = s.articleRepository.CreateEsIndex(ctx, index, esMapper)
		if err != nil {
			return 0, v1.ErrCreateEsIndexFailed
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
	// 判断团队下是否有知识库
	_, total, err := s.kbRepository.GetKBListByTeamId(ctx, teamId, 1, 200)
	if err != nil {
		return v1.ErrGetCategoryListFailed
	}
	if total > 0 {
		// 删除知识库
		if err := s.kbRepository.DeleteKBByTeamID(ctx, teamId); err != nil {
			return v1.ErrDeleteKnowledgeFailed
		}
	}
	// 删除es索引团队
	index := enums.Team_knowledge_index + strconv.Itoa(int(teamId))
	err = s.articleRepository.DeleteEsIndex(ctx, index)
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
	// 根据用户ID查询团队成员列表
	memberList, total, err := s.teamRepository.GetTeamMemberListByUserID(ctx, userId, index, size)
	if err != nil {
		return nil, v1.ErrGetTeamMemberListFailed
	}
	// 构建团队ID列表
	var teamIds []uint
	for _, member := range memberList {
		teamIds = append(teamIds, member.TeamID)
	}
	// 查询团队列表，根据size切片查询
	teams, err := s.teamRepository.GetTeamListByIDs(ctx, teamIds)
	if err != nil {
		return nil, v1.ErrGetTeamListFailed
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
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 发生其他错误（不是"未找到记录"的错误），直接返回
		return err
	}
	if err == nil {
		// 查询成功，说明用户已在团队中
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
	// 判断用户是否为团队负责人或者管理员
	if member.Role == enums.LEADER || member.Role == enums.ADMIN {
		// 判断团队下是否有其他管理员
		members, err := s.teamRepository.GetTeamMemberListByTeamID(ctx, teamId)
		if err != nil {
			return v1.ErrGetTeamMemberListFailed
		}
		var adminCount int
		for _, member := range members {
			if member.Role == enums.ADMIN || member.Role == enums.LEADER {
				adminCount++
			}
		}
		if adminCount <= 1 {
			return v1.ErrAtLeastOneAdmin
		}
	}
	// 删除成员
	err = s.teamRepository.DeleteTeamMember(ctx, member.MemberID)
	if err != nil {
		return v1.ErrDeleteMemberFailed
	}
	return nil
}

func (s *teamService) ApplyJoinTeam(ctx *gin.Context, userID string, req *v1.ApplyJoinTeamReq) (uint, error) {
	// 判断用户是否在团队中
	_, err := s.teamRepository.GetMemberByTeamIDAndUserID(ctx, req.TeamID, userID)
	if err == nil {
		return 0, v1.ErrMemberExist
	}
	// 判断用户是否已经提交过申请
	_, err = s.teamRepository.GetApplyByTeamIDAndUserID(ctx, req.TeamID, userID)
	if err == nil {
		return 0, v1.ErrApplyExisted
	}
	// 构建申请表
	apply := &model.TeamApplications{
		TeamID:      req.TeamID,
		ApplicantID: userID,
		Status:      enums.WAITING,
		Reason:      &req.Reason,
	}
	applyId, err := s.teamRepository.CreateTeamApply(ctx, apply)
	if err != nil {
		return 0, v1.ErrApplyFailed
	}
	return applyId, nil
}
