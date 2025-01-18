package model

import (
	"gorm.io/gorm"
	v1 "projectName/api/v1"
	"time"
)

// Article 文章结构体，GORM 数据模型
type Article struct {
	ArticleID       uint            `gorm:"primaryKey;autoIncrement"`   // 文章的唯一ID，数据库主键
	Title           string          `gorm:"type:varchar(255);not null"` // 文章标题
	Content         string          `gorm:"type:text;not null"`         // 文章内容
	ContentShort    string          `gorm:"type:varchar(255)"`          // 文章摘要
	UserID          string          `gorm:"type:varchar(255);not null"` // 用户ID
	CategoryID      uint            `gorm:"unique;not null"`            // 分类ID
	Importance      int             `gorm:"type:int;default:0"`         // 文章重要性
	VisibleRange    string          `gorm:"type:varchar(255);not null"` // 可见范围
	CommentDisabled bool            `gorm:"type:boolean;default:false"` // 是否禁用评论
	SourceURI       string          `gorm:"type:varchar(255)"`          // 文章外链
	Status          int             `gorm:"type:int;default:0"`         // 文章状态
	UploadedFiles   []v1.FileUpload `gorm:"type:json"`                  // 上传的文件列表
	CreatedAt       time.Time       `gorm:"autoCreateTime" `            // 文章创建时间
	UpdatedAt       time.Time       `gorm:"autoUpdateTime" `            // 文章更新时间
	DeletedAt       gorm.DeletedAt  `gorm:"index"`
}

func (m *Article) TableName() string {
	return "kb_article"
}
