package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id        uint   `gorm:"primarykey"`
	UserId    int64  `gorm:"unique;not null"`
	Phone     string `gorm:"not null"`
	Nickname  string `gorm:"not null"`
	Password  string `gorm:"not null"`
	RoleType  int    `gorm:"not null"` // 0: 普通用户，1: 学校用户，2: 学校管理员 3: 超级管理员
	Email     string
	CollegeId uint
	StudentId string
	IsDeleted int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) TableName() string {
	return "users"
}
