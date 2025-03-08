package vo

type CategoryView struct {
	CId          uint           `json:"cid" gorm:"column:category_id"`            // 分类ID
	CategoryName string         `json:"categoryName" gorm:"column:category_name"` // 分类名
	ParentId     uint           `json:"parentId" gorm:"column:parent_id"`         // 父分类ID
	KbID         uint           `json:"kbId" gorm:"column:kb_id"`                 // 知识库ID
	Level        int            `json:"level" gorm:"column:level"`                // 层级
	Children     []CategoryView `json:"children" gorm:"-"`                        // 子分类
}
