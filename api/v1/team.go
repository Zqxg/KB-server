package v1

import (
	"time"
)

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
type TeamData struct {
	TeamID      uint      `json:"team_id"`      // 团队ID
	TeamName    string    `json:"team_name"`    // 团队名
	Description string    `json:"description"`  // 描述
	CreatedBy   string    `json:"created_by"`   // 创建者ID
	CreatorName string    `json:"creator_name"` // 创建者昵称
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`   // 更新时间
}

type MemberData struct {
	MemberID   uint      `json:"member_id"`   // 成员ID
	TeamID     uint      `json:"team_id"`     // 团队ID
	UserID     string    `json:"member_name"` // 成员ID
	NickName   string    `json:"nick_name"`   // 成员昵称
	RoleType   string    `json:"role_type"`   // 角色类型
	JoinTime   time.Time `json:"join_time"`   // 加入时间
	UpdateTime time.Time `json:"update_time"` // 更新时间
}

type GetTeamListResp struct {
	TeamList []TeamData `json:"team_list"`
	PageResponse
}

type GetTeamInfoResp struct {
	Team   TeamData      `json:"team"`
	Member []*MemberData `json:"member"`
}

type GetUserTeamListResp struct {
	TeamList []TeamData `json:"team_list"`
	PageResponse
}

type GetTeamMemberListResp struct {
	MemberList []*MemberData `json:"member_list"`
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

type ApplyJoinTeamReq struct {
	TeamID uint   `json:"team_id"` // 团队ID
	Reason string `json:"reason"`  // 申请理由
}
type ApplyJoinTeamResp struct {
	ApplyID uint `json:"apply_id"` // 申请ID
}

type GetTeamApplyListReq struct {
	TeamID     uint `json:"team_id"`     // 团队ID
	Status     int  `json:"status"`      // 申请状态
	IsPersonal bool `json:"is_personal"` // 是否是个人申请
	PageRequest
}
type TeamApplyData struct {
	ApplyID       uint      `json:"apply_id"`       // 申请ID
	TeamID        uint      `json:"team_id"`        // 团队ID
	TeamName      string    `json:"team_name"`      // 团队名
	ApplicantID   string    `json:"applicant_id"`   // 申请人ID
	ApplicantName string    `json:"applicant_name"` // 用户昵称
	ReviewerID    string    `json:"reviewer_id"`    // 审核人ID
	ReviewerName  string    `json:"reviewer_name"`  // 审核人昵称
	Reason        string    `json:"reason"`         // 申请理由
	Status        int       `json:"status"`         // 申请状态
	CreateTime    time.Time `json:"create_time"`    // 创建时间
	UpdateTime    time.Time `json:"update_time"`    // 更新时间
}

type GetTeamApplyListResp struct {
	ApplyList []*TeamApplyData `json:"apply_list"`
	PageResponse
}

type HandleTeamApplyReq struct {
	ApplyID uint `json:"apply_id"` // 申请ID
	TeamID  uint `json:"team_id"`  // 团队ID
	Status  int  `json:"status"`   // 申请状态
}
