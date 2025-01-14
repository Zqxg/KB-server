package enums

// 文章Status
const (
	StatusPublished     = 1 // 已发布
	StatusDeleted       = 2 // 已删除
	StatusPendingReview = 3 // 待审核
	StatusRejected      = 4 // 已驳回
	StatusScheduled     = 5 // 已计划 设置了定时发布的时间
)
