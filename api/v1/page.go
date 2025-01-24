package v1

// PageResponse 通用分页数据结构体
type PageResponse struct {
	TotalCount int64 `json:"totalCount"` // 总记录数
	PageIndex  int   `json:"pageIndex"`  // 当前页码
	PageSize   int   `json:"pageSize"`   // 每页大小
}

type PageRequest struct {
	PageIndex int `form:"pageIndex" json:"pageIndex"` // 当前页码
	PageSize  int `form:"pageSize" json:"pageSize"`   // 每页大小
}
