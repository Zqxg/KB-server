package model

import "time"

type UserAuth struct {
	Id          uint       `gorm:"primaryKey"`
	UserId      string     `gorm:"not null;index"`                     // 用户ID，外键关联sys_users表
	RequestType int        `gorm:"not null"`                           // 认证类型（1: 学生认证, 2: 管理员认证）
	Status      int        `gorm:"not null;default:0"`                 // 认证请求状态（0: 待处理, 1: 已批准, 2: 已拒绝, 3: 认证失败）
	ApplyTime   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"` // 申请时间
	DisposeTime *time.Time `gorm:"default:null"`                       // 处理时间
	AdminId     *string    `gorm:"default:null"`                       // 处理该请求的管理员ID
	Remarks     *string    `gorm:"type:text;default:null"`             // 备注，允许为空
	CollegeId   *uint      `gorm:"default:null"`                       // 学校ID，允许为空
	StudentId   *string    `gorm:"default:null"`                       // 学生ID，允许为空
}

func (ua *UserAuth) TableName() string {
	return "sys_user_auths"
}
