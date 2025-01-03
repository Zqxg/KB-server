package model

import (
	"time"
)

type College struct {
	Id          uint       `gorm:"primarykey"`                // 学院ID
	CollegeId   uint       `gorm:"unique;not null"`           // 学院唯一标识ID
	CollegeName string     `gorm:"not null"`                  // 学院名称
	Description string     `gorm:"type:text"`                 // 学院描述
	IsDeleted   int        `gorm:"default:0"`                 // 是否删除（软删除字段）
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP"` // 创建时间
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP"` // 更新时间
	DeletedAt   *time.Time `gorm:"index"`                     // 软删除时间
}

func (m *College) TableName() string {
	return "sys_colleges"
}
