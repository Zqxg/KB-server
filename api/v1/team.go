package v1

import "projectName/internal/model"

type CreateTeamRequest struct {
	TeamName    string `json:"team_name"`   // 团队名
	Description string `json:"description"` // 描述
}
type CreateTeamResp struct {
	TeamID uint `json:"team_id"` // 团队ID
}

type UpdateTeamRequest struct {
	TeamID      uint   `json:"team_id"`     // 团队ID
	TeamName    string `json:"team_name"`   // 团队名
	Description string `json:"description"` // 描述
}

type DeleteTeamRequest struct {
	TeamID uint `json:"team_id"` // 团队ID
}

type GetTeamListReq struct {
	TeamName  string `json:"team_name"`  // 团队名
	CreatedBy string `json:"created_by"` // 创建者ID
	PageRequest
}

type GetTeamListResp struct {
	TeamList []model.Team `json:"team_list"`
	PageResponse
}

type GetTeamInfoResp struct {
	Team   model.Team      `json:"team"`
	Member []*model.Member `json:"member"`
}

type GetUserTeamListResp struct {
	TeamList []model.Team `json:"team_list"`
	PageResponse
}

type GetTeamMemberListResp struct {
	MemberList []*model.Member `json:"member_list"`
}

type AddTeamMemberReq struct {
	TeamID   uint   `json:"team_id"`   // 团队ID
	MemberID string `json:"member_id"` // 成员ID
	RoleType string `json:"role_type"` // 角色类型
}

type DeleteTeamMemberReq struct {
	TeamID   uint   `json:"team_id"`   // 团队ID
	MemberID string `json:"member_id"` // 成员ID
}

type UpdateTeamMemberRoleReq struct {
	TeamID   uint   `json:"team_id"`   // 团队ID
	MemberID string `json:"member_id"` // 成员ID
	RoleType string `json:"role_type"` // 角色类型
}

type QuitTeamReq struct {
	TeamID uint `json:"team_id"` // 团队ID
}
