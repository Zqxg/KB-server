package v1

// PageResponse 通用分页数据结构体
type PageResponse struct {
	TotalCount int64 `json:"total_count"` // 总记录数
	PageIndex  int   `json:"page_index"`  // 当前页码
	PageSize   int   `json:"page_size"`   // 每页大小
}

type PageRequest struct {
	PageIndex int `form:"pageIndex" json:"page_index"` // 当前页码
	PageSize  int `form:"pageSize" json:"page_size"`   // 每页大小
}
