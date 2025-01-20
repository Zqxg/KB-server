package v1

import (
	"projectName/internal/model/vo"
)

// CreateArticleRequest 用于接收创建文章请求的数据
type CreateArticleRequest struct {
	Title           string       `json:"title" binding:"required"`        // 文章标题
	Content         string       `json:"content" binding:"required"`      // 文章内容
	ContentShort    string       `json:"contentShort"`                    // 文章摘要
	AuthorID        string       `json:"authorId" binding:"required"`     // 作者ID
	CategoryID      uint         `json:"categoryId"`                      // 文章分类ID
	Importance      int          `json:"importance"`                      // 文章重要性
	VisibleRange    string       `json:"visibleRange" binding:"required"` // 可见范围
	CommentDisabled bool         `json:"commentDisabled"`                 // 是否禁用评论
	SourceURI       string       `json:"sourceUri"`                       // 文章外链
	UploadedFiles   []FileUpload `json:"uploadedFiles"`                   // 上传的文件列表
}

// FileUpload 用于接收上传文件的信息
type FileUpload struct {
	FileName string `json:"fileName" ` // 文件名
	FileURL  string `json:"fileUrl" `  // 文件URL
}

type CreateArticleResponseData struct {
	ArticleID int `json:"articleId"` // 文章ID
}

type ArticleResponseData struct {
	ArticleID       uint         `json:"articleId"`       //文章id
	Title           string       `json:"title" `          // 文章标题
	Content         string       `json:"content" `        // 文章内容
	ContentShort    string       `json:"contentShort"`    // 文章摘要
	Author          string       `json:"author" `         // 作者
	Category        string       `json:"category"`        // 文章分类
	Importance      int          `json:"importance"`      // 文章重要性
	VisibleRange    string       `json:"visibleRange" `   // 可见范围
	CommentDisabled bool         `json:"commentDisabled"` // 是否禁用评论
	SourceURI       string       `json:"sourceUri"`       // 文章外链
	UploadedFiles   []FileUpload `json:"uploadedFiles"`   // 上传的文件列表
	CreatedAt       string       `json:"createdAt"`       // 文章创建时间
	UpdatedAt       string       `json:"updateAt"`        // 文章更新时间
	//Tags            []Tags       `json:"tags"`            //todo：文章标签
}

type CategoryList []vo.CategoryView
type CategoryData struct {
	CategoryList
}

type GetArticleRequest struct {
	ArticleID int `json:"articleId"` // 文章ID
}
