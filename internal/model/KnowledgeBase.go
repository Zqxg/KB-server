package model

import (
	"gorm.io/gorm"
	"time"
)

// KnowledgeBase 知识库表
type KnowledgeBase struct {
	KbID      uint           `gorm:"primaryKey;autoIncrement"`   // 知识库ID
	KbName    string         `gorm:"type:varchar(255);not null"` // 知识库名称
	TeamID    *uint          `gorm:""`                           // 团队ID (NULL 表示公共知识库)
	UserID    *string        `gorm:"type:varchar(255)"`          // 用户ID (NULL 表示非私人知识库)
	IsPublic  bool           `gorm:"type:boolean;default:false"` // 是否为公共知识库
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (KnowledgeBase) TableName() string {
	return "kb_knowledgeBase"
}
