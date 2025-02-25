package v1

import (
	"projectName/internal/model"
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

type ArticleData struct {
	ArticleID       uint         `json:"articleId"`       //文章id
	Title           string       `json:"title" `          // 文章标题
	Content         string       `json:"content" `        // 文章内容
	ContentShort    string       `json:"contentShort"`    // 文章摘要
	Author          string       `json:"author" `         // 作者
	Category        string       `json:"category"`        // 文章分类
	CategoryID      uint         `json:"categoryId"`      // 文章分类ID
	Importance      int          `json:"importance"`      // 文章重要性
	VisibleRange    string       `json:"visibleRange" `   // 可见范围
	CommentDisabled bool         `json:"commentDisabled"` // 是否禁用评论
	SourceURI       string       `json:"sourceUri"`       // 文章外链
	UploadedFiles   []FileUpload `json:"uploadedFiles"`   // 上传的文件列表
	Status          int          `json:"status"`          // 文章状态
	CreatedAt       string       `json:"createdAt"`       // 文章创建时间
	UpdatedAt       string       `json:"updateAt"`        // 文章更新时间
	//Tags            []Tags       `json:"tags"`            //todo：文章标签
}

type CategoryList []vo.CategoryView
type CategoryData struct {
	CategoryList
}

type GetArticleRequest struct {
	ArticleID uint `json:"articleId"` // 文章ID
}

type UpdateArticleRequest struct {
	ArticleID uint `json:"articleId"` // 文章ID
	CreateArticleRequest
}

type DelArticleListReq struct {
	ArticleIDList []uint `json:"articleIDList"` // 文章ID列表
}

type DeleteArticleResponseData struct {
	DeletedCount int `json:"deletedCount"` // 删除的文章数量
}

type DeleteArticleRequest struct {
	ArticleID uint `json:"articleId"` // 文章ID
}

type GetArticleListByCategoryReq struct {
	CategoryID uint `json:"categoryId"` // 文章分类ID
	PageRequest
}

type ArticleList struct {
	ArticleDataList []*ArticleData
	PageResponse
}

type GetUserArticleListReq struct {
	Title      string `json:"title"`      // 文章标题
	CategoryID uint   `json:"categoryId"` // 文章分类ID
	CreatedAt  string `json:"createdAt"`  // 文章创建时间
	CreatedEnd string `json:"CreatedEnd"` // 文章结束时间
	Status     int    `json:"status"`     // 文章状态
	PageRequest
}

type GetArticleListByEsReq struct {
	PageRequest
	Content     string   `json:"content"`     // 搜索的内容关键词
	Title       string   `json:"title"`       // 搜索的标题
	Keywords    []string `json:"keywords"`    // 搜索的关键字数组
	PhraseMatch bool     `json:"phraseMatch"` // 是否启用短语匹配
	AdvSearch   bool     `json:"advSearch"`   // 是否启用高级搜索
	Column      string   `json:"column"`      // 排序字段，通常是 "_score"
	Order       string   `json:"order"`       // 排序方式，"asc" 或 "desc"
	Categories  []int    `json:"categories"`  // 分类id，用于筛选
	SearchMode  string   `json:"searchMode"`  // 搜索模式标识（例如 "0" 代表普通搜索 "1"）
}

type SearchArticleResp struct {
	PageResponse
	Articles []ArticleSearchInfo `json:"articles"` // 文章列表
}

type ArticleSearchInfo struct {
	EsArticle model.EsArticle // 文章数据
	Score     float64         `json:"score"` // 评分（例如：基于ES的相关度评分）
	//Tags            []string     `json:"tags"`             // 文章标签（可选）
	//Comments        *int         `json:"comments"`         // 评论数（可选）
	//Views           *int         `json:"views"`            // 浏览量（可选）
}
