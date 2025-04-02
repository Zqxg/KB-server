package v1

import (
	"time"
)

// CreateArticleRequest 用于接收创建文章请求的数据
type CreateArticleRequest struct {
	Title           string       `json:"title" binding:"required"`     // 文章标题
	Content         string       `json:"content" binding:"required"`   // 文章内容
	ContentShort    string       `json:"content_short"`                // 文章摘要
	AuthorID        string       `json:"author_id" binding:"required"` // 作者ID
	KBID            uint         `json:"kb_id"`                        // 知识库ID
	CategoryID      uint         `json:"category_id"`                  // 文章分类ID
	Importance      int          `json:"importance"`                   // 文章重要性
	Status          int          `json:"status" binding:"required"`    // 文章状态
	CommentDisabled bool         `json:"comment_disabled"`             // 是否禁用评论
	SourceURI       string       `json:"source_uri"`                   // 文章外链
	UploadedFiles   []FileUpload `json:"uploaded_files"`               // 上传的文件列表
}

// FileUpload 用于接收上传文件的信息
type FileUpload struct {
	FileName string `json:"file_name" ` // 文件名
	FileURL  string `json:"file_url" `  // 文件URL
}

type CreateArticleResponseData struct {
	ArticleID int `json:"article_id"` // 文章ID
}

type ArticleData struct {
	ArticleID       uint         `json:"article_id"`       //文章id
	Title           string       `json:"title" `           // 文章标题
	Content         string       `json:"content" `         // 文章内容
	ContentShort    string       `json:"content_short"`    // 文章摘要
	Author          string       `json:"author" `          // 作者
	Category        string       `json:"category"`         // 文章分类
	CategoryID      uint         `json:"category_id"`      // 文章分类ID
	KBID            uint         `json:"kb_id"`            // 知识库ID
	KBName          string       `json:"kb_name"`          // 知识库名称
	Importance      int          `json:"importance"`       // 文章重要性
	CommentDisabled bool         `json:"comment_disabled"` // 是否禁用评论
	SourceURI       string       `json:"source_uri"`       // 文章外链
	UploadedFiles   []FileUpload `json:"uploaded_files"`   // 上传的文件列表
	Status          int          `json:"status"`           // 文章状态
	CreatedAt       string       `json:"created_at"`       // 文章创建时间
	UpdatedAt       string       `json:"update_at"`        // 文章更新时间
	//Tags            []Tags       `json:"tags"`            //todo：文章标签
}

type GetArticleRequest struct {
	ArticleID uint `json:"article_id"` // 文章ID
}

type UpdateArticleRequest struct {
	ArticleID uint `json:"article_id"` // 文章ID
	CreateArticleRequest
}

type DelArticleListReq struct {
	ArticleIDList []uint `json:"article_id_List"` // 文章ID列表
}

type DeleteArticleResponseData struct {
	DeletedCount int `json:"deleted_count"` // 删除的文章数量
}

type DeleteArticleRequest struct {
	ArticleID uint `json:"article_id"` // 文章ID
}

type GetArticleListByCategoryReq struct {
	CategoryID uint `json:"article_id"` // 文章分类ID
	PageRequest
}

type ArticleList struct {
	ArticleDataList []*ArticleData `json:"article_data_list"`
	PageResponse
}

type GetUserArticleListReq struct {
	Title      string `json:"title"`       // 文章标题
	CategoryID uint   `json:"category_id"` // 文章分类ID
	KBID       uint   `json:"kb_id"`       // 知识库ID
	CreatedAt  string `json:"created_at"`  // 文章创建时间
	CreatedEnd string `json:"created_end"` // 文章结束时间
	Status     int    `json:"status"`      // 文章状态
	PageRequest
}

type GetArticleListByEsReq struct {
	PageRequest
	Content         string   `json:"content"`           // 搜索的内容关键词
	Title           string   `json:"title"`             // 搜索的标题
	Keywords        []string `json:"keywords"`          // 搜索的关键字数组
	PhraseMatch     bool     `json:"phrase_match"`      // 是否启用短语匹配
	AdvSearch       bool     `json:"adv_search"`        // 是否启用高级搜索
	Column          string   `json:"column"`            // 排序字段，通常是 "_score"
	Order           string   `json:"order"`             // 排序方式，"asc" 或 "desc"
	Importance      string   `json:"importance"`        // 文章重要性
	CreateTimeStart string   `json:"create_time_start"` // 文章创建时间
	CreateTimeEnd   string   `json:"create_time_end"`   // 文章结束时间
	Categories      []int    `json:"categories"`        // 分类id，用于筛选
}

type SearchArticleResp struct {
	PageResponse
	Articles []ArticleSearchInfo `json:"articles"` // 文章列表
}

type ArticleSearchInfo struct {
	ArticleID       uint      `json:"article_id"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	ContentShort    string    `json:"content_short"`
	Author          string    `json:"author"`
	Category        string    `json:"category"`
	KBName          string    `json:"kb_name"`
	TeamName        string    `json:"team_name"`
	Importance      int       `json:"importance"`
	CommentDisabled bool      `json:"comment_disabled"`
	SourceURI       string    `json:"source_uri"`
	Status          int       `json:"status"`
	UploadedFile    bool      `json:"uploaded_file"`
	CreatedAt       time.Time `json:"created_at"` // 使用 sql.NullTime
	UpdatedAt       time.Time `json:"updated_at"` // 使用 sql.NullTime
	Score           float64   `json:"score"`      // 评分（例如：基于ES的相关度评分）
	//Tags            []string     `json:"tags"`             // 文章标签（可选）
	//Comments        *int         `json:"comments"`         // 评论数（可选）
	//Views           *int         `json:"views"`            // 浏览量（可选）
}
