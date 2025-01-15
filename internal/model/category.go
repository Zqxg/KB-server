package model

import (
	"gorm.io/gorm"
	"time"
)

type Category struct {
	CId          uint   `gorm:"primaryKey;autoIncrement"`
	CategoryName string `gorm:"type:varchar(255);not null"`
	ParentId     uint   `gorm:"type:int;default:0"`
	IsDeleted    int    `gorm:"default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (m *Category) TableName() string {
	return "kb_category"
}
