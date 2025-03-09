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
