package enums

// memberRole 成员角色
const (
	LEADER = "leader" // 负责人
	ADMIN  = "admin"  // 管理员
	MEMBER = "member" // 成员
)

// teamStatus 团队状态
const (
	StatusPending   = 0 // 待审核
	StatusWithdrawn = 1 // 已撤回
	StatusApproved  = 2 // 已通过
	StatusRejected  = 3 // 已拒绝
)
