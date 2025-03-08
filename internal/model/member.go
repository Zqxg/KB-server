package model

import (
	"gorm.io/gorm"
	"time"
)

// Member 团队成员表
type Member struct {
	MemberID  uint           `gorm:"primaryKey;autoIncrement"`                      // 成员ID
	TeamID    uint           `gorm:"not null"`                                      // 所属团队ID
	UserID    string         `gorm:"type:varchar(255);not null"`                    // 成员用户ID
	Role      string         `gorm:"type:enum('leader','admin','member');not null"` // 成员角色
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Member) TableName() string {
	return "sys_member"
}
