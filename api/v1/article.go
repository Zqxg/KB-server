package v1

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
	FileName string `json:"fileName" binding:"required"` // 文件名
	FileURL  string `json:"fileUrl" binding:"required"`  // 文件URL
}

type CreateArticleResponseData struct {
	ArticleID uint `json:"articleId"` // 文章ID
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
	//Tags            []Tags       `json:"tags"`            //todo：文章标签
}
