package enums

// 文章Status
const (
	StatusDraft         = 1  // 草稿
	StatusPublished     = 2  // 已发布
	StatusPendingReview = 3  // 待审核
	StatusDeleted       = 99 // 已删除
	StatusAll           = -1
)

// GetStatus 判断是否存在 存在返回当前状态 不存在返回默认状态
func GetStatus(status int) int {
	if IsStatusExist(status) {
		return status
	}
	return StatusAll
}

func IsStatusExist(status int) bool {
	switch status {
	case StatusDraft, StatusPublished, StatusPendingReview, StatusDeleted:
		return true
	default:
		return false
	}
}
