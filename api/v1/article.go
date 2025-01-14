package v1

// GetArticleRequest 用于接收创建文章请求的数据
type GetArticleRequest struct {
	Title           string       `json:"title" binding:"required"`        // 文章标题
	Content         string       `json:"content" binding:"required"`      // 文章内容
	ContentShort    string       `json:"contentShort"`                    // 文章摘要
	AuthorID        string       `json:"author" binding:"required"`       // 作者ID
	CategoryID      uint         `json:"categoryID"`                      // 文章分类ID
	Importance      int          `json:"importance"`                      // 文章重要性
	VisibleRange    string       `json:"visibleRange" binding:"required"` // 可见范围
	CommentDisabled bool         `json:"commentDisabled"`                 // 是否禁用评论
	SourceURI       string       `json:"sourceUri"`                       // 文章外链
	UploadedFiles   []FileUpload `json:"uploadedFiles"`                   // 上传的文件列表
}

// FileUpload 用于接收上传文件的信息
type FileUpload struct {
	FileName string `json:"fileName" binding:"required"` // 文件名
	FileURL  string `json:"fileUrl" binding:"required"`  // 文件URL
}

type GetArticleResponseData struct {
	ArticleID string `json:"articleId"` // 文章ID
}
