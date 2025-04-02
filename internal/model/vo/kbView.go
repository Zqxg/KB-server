package vo

import "time"

type KbKnowledgeBaseView struct {
	KBID      uint      `json:"kb_id" gorm:"column:kb_id"`
	KbName    string    `json:"kb_name" gorm:"column:kb_name"`
	TeamID    uint      `json:"team_id,omitempty" gorm:"column:team_id"` // 团队ID，可能为空
	UserID    string    `json:"user_id,omitempty" gorm:"column:user_id"` // 用户ID，可能为空（私人知识库）
	IsPublic  bool      `json:"is_public" gorm:"column:is_public"`
	CreatedBy string    `json:"created_by" gorm:"column:created_by"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	TeamName  string    `json:"team_name,omitempty" gorm:"column:team_name"` // 团队名称
	UserName  string    `json:"user_name,omitempty" gorm:"column:user_name"` // 用户昵称
	KBType    string    `json:"kb_type" gorm:"column:kb_type"`               // 知识库类型：私人/公共/团队
}

// 指定表名（视图名）
func (KbKnowledgeBaseView) TableName() string {
	return "kb_knowledge_base_view"
}
