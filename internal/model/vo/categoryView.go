package vo

type CategoryView struct {
	CId          uint           `json:"cid" gorm:"column:category_id"`
	CategoryName string         `json:"categoryName" gorm:"column:category_name"`
	ParentId     uint           `json:"parentId" gorm:"column:parent_id"`
	Level        int            `json:"level" gorm:"column:level"`
	Children     []CategoryView `json:"children" gorm:"-"`
}
