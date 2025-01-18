package v1

// PageData 通用分页数据结构体
type PageData struct {
	Data       interface{} `json:"data"`       // 当前页的数据
	TotalCount int         `json:"totalCount"` // 总记录数
	TotalPage  int         `json:"totalPage"`  // 总页数
	PageIndex  int         `json:"pageIndex"`  // 当前页码
	PageSize   int         `json:"pageSize"`   // 每页大小
}
