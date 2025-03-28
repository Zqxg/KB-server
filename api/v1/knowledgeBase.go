package v1

import (
	"projectName/internal/model/vo"
	"time"
)

type CreateKBRequest struct {
	TeamID uint   `json:"team_id" binding:"required"`
	Name   string `json:"name" binding:"required"`
}

type CreateKBResp struct {
	KBID uint `json:"kb_id"`
}

type UpdateKBNameReq struct {
	KBID uint   `json:"kb_id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type DeleteKBReq struct {
	KBID uint `json:"kb_id" binding:"required"`
}

type GetKBInfoResp struct {
	KBID      uint      `json:"kb_id"`
	KbName    string    `json:"kb_name"`
	TeamID    uint      `json:"team_id"`
	UserID    string    `json:"user_id"`
	TeamName  string    `json:"team_name"`
	UserName  string    `json:"user_name"`
	IsPublic  bool      `json:"is_public"`
	KBType    string    `json:"kb_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type GetKBListByTeamIdReq struct {
	TeamID uint `json:"team_id" binding:"required"`
	PageRequest
}

type GetKBListByTeamIdResp struct {
	KBList []GetKBInfoResp `json:"kb_list"`
	PageResponse
}

type CreateCategoryReq struct {
	KBID         uint   `json:"kb_id" binding:"required"`
	CategoryName string `json:"category_name" binding:"required"`
	ParentId     uint   `json:"parent_id"`
}

type UpdateCategoryReq struct {
	CID          uint   `json:"cid" binding:"required"`
	CategoryName string `json:"category_name" binding:"required"`
}

type DeleteCategoryReq struct {
	CID uint `json:"cid" binding:"required"`
}

type CategoryList []vo.CategoryView
type CategoryData struct {
	CategoryList
}

type GetKBListByTypeReq struct {
	KBType string `json:"kb_type"`
}

type KBList struct {
	KBList []*vo.KbKnowledgeBaseView `json:"kb_list"`
}
