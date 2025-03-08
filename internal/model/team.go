package model

import (
	"gorm.io/gorm"
	"time"
)

// Team 团队表
type Team struct {
	TeamID      uint           `gorm:"primaryKey;autoIncrement"`   // 团队ID
	TeamName    string         `gorm:"type:varchar(255);not null"` // 团队名
	Description string         `gorm:"type:text"`                  // 描述
	CreatedBy   string         `gorm:"type:varchar(255);not null"` // 创建者ID
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Team) TableName() string {
	return "sys_team"
}
