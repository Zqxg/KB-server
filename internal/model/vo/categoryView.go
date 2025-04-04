package vo

type CategoryView struct {
	CId          uint           `json:"cid" gorm:"column:category_id"`             // 分类ID
	CategoryName string         `json:"category_name" gorm:"column:category_name"` // 分类名
	ParentId     uint           `json:"parent_id" gorm:"column:parent_id"`         // 父分类ID
	KbID         uint           `json:"kb_id" gorm:"column:kb_id"`                 // 知识库ID
	Level        int            `json:"level" gorm:"column:level"`                 // 层级
	Children     []CategoryView `json:"children" gorm:"-"`                         // 子分类
}

func (CategoryView) TableName() string {
	return "kb_category_view"
}
