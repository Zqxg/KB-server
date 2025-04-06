package model

import (
	"gorm.io/gorm"
	"time"
)

type TeamApplications struct {
	ApplicationID uint           `gorm:"primaryKey;autoIncrement"`
	TeamID        uint           `gorm:"type:bigint unsigned;not null"`
	ApplicantID   string         `gorm:"type:varchar(255);collate:utf8mb4_unicode_ci;not null"`
	ReviewerID    *string        `gorm:"type:varchar(255);collate:utf8mb4_unicode_ci"`
	Status        int            `gorm:"type:varchar(255);collate:utf8mb4_unicode_ci;not null"`
	Reason        *string        `gorm:"type:varchar(255);collate:utf8mb4_unicode_ci"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (TeamApplications) TableName() string {
	return "sys_team_applications"
}
