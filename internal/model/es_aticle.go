package model

import (
	"time"
)

// EsArticle 文章结构
type EsArticle struct {
	ArticleID       uint      `json:"article_id"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	ContentShort    string    `json:"content_short"`
	UserID          string    `json:"user_id"`
	CategoryID      uint      `json:"category_id"`
	Importance      int       `json:"importance"`
	VisibleRange    string    `json:"visible_range"`
	CommentDisabled bool      `json:"comment_disabled"`
	SourceURI       string    `json:"source_uri"`
	Status          int       `json:"status"`
	UploadedFiles   string    `json:"uploaded_files"`
	CreatedAt       time.Time `json:"created_at"` // 使用 sql.NullTime
	UpdatedAt       time.Time `json:"updated_at"` // 使用 sql.NullTime
}
