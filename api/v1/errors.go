package v1

var (
	// common errors
	ErrSuccess             = newError(0, "操作成功")
	ErrBadRequest          = newError(400, "请求参数错误")
	ErrUnauthorized        = newError(401, "未授权")
	ErrPermissionDenied    = newError(403, "权限不足")
	ErrNotFound            = newError(404, "未找到")
	ErrInternalServerError = newError(500, "内部服务器错误")

	// more biz errors
	//ErrEmailAlreadyUse = newError(1001, "邮箱已存在")
	ErrPhoneAlreadyUse = newError(1002, "手机号已存在")
	ErrPhoneFormat     = newError(1003, "手机号格式错误")
	ErrPasswordFormat  = newError(1004, "密码格式错误")
	ErrDecryptPassword = newError(1005, "错误解密密码")
	ErrGetTokenFail    = newError(1006, "获取token失败")
	ErrUserNotExist    = newError(1007, "用户不存在")
	ErrLogoutFail      = newError(1008, "用户退出失败")
	ErrCancelFail      = newError(1009, "用户注销失败")
	ErrSearchFailed    = newError(1010, "搜索失败")
	//ErrEmailFormat     = newError(1010, "邮箱格式错误")

	ErrArticleNotExist     = newError(1101, "文章不存在")
	ErrUpdateArticleFailed = newError(1102, "修改文章失败")
	ErrArticleStatusError  = newError(1103, "文章状态异常")

	ErrUploadFileFailed      = newError(1104, "上传文件序列化失败")
	ErrDeserializeFileFailed = newError(1105, "上传文件反序列化失败")

	// 2000 错误码
	ErrInvalidCaptcha = newError(2000, "验证码错误")

	// 2100 知识库
	ErrKnowledgeNotExist = newError(2100, "知识库不存在")
	ErrKnowledgeExist    = newError(2101, "知识库已存在")
	//ErrKnowledgeEmpty    = newError(2102, "知识库为空")
	// 知识库删除错误
	ErrDeleteKnowledgeFailed = newError(2103, "删除知识库失败")
	//ErrDeleteKnowledgeExist  = newError(2104, "知识库存在文章")
	//ErrDeleteKnowledgePublic = newError(2105, "公共知识库不能删除")
	// 知识库创建错误
	ErrCreateKnowledgeFailed = newError(2106, "创建知识库失败")
	//ErrCreateKnowledgeExist  = newError(2107, "知识库已存在")
	//ErrCreateKnowledgePublic = newError(2108, "公共知识库不能创建")
	// 知识库更新错误
	ErrUpdateKnowledgeFailed = newError(2109, "更新知识库失败")
	//ErrUpdateKnowledgePublic = newError(2110, "公共知识库不能更新")
	ErrNotPublicKnowledge = newError(2110, "非公共知识库")
	// 2111 知识库分类
	ErrCategoryNotExist      = newError(2111, "分类不存在")
	ErrUpdateCategoryFailed  = newError(2112, "更新分类失败")
	ErrCreateCategoryFailed  = newError(2113, "创建分类失败")
	ErrDeleteCategoryFailed  = newError(2114, "删除分类失败")
	ErrGetCategoryListFailed = newError(2115, "获取分类列表失败")
	ErrCategoryNotMatchKB    = newError(2116, "分类不属于该知识库")

	// 2200 team
	ErrTeamNotExist = newError(2200, "团队不存在")
	//ErrTeamExist         = newError(2201, "团队已存在")
	//ErrTeamEmpty         = newError(2202, "团队为空")
	ErrDeleteTeamFailed  = newError(2203, "删除团队失败")
	ErrCreateTeamFailed  = newError(2206, "创建团队失败")
	ErrUpdateTeamFailed  = newError(2209, "更新团队失败")
	ErrGetTeamInfoFailed = newError(2210, "获取团队信息失败")
	ErrGetTeamListFailed = newError(2211, "获取团队列表失败")
	ErrTeamNameTooShort  = newError(2212, "团队名称过短")
	// 2300 成员
	ErrNoTeamAdminPermission   = newError(2304, "无团队管理员权限")
	ErrMemberNotExist          = newError(2300, "成员不存在")
	ErrMemberExist             = newError(2301, "成员已存在")
	ErrDeleteMemberFailed      = newError(2303, "删除成员失败")
	ErrUpdateMemberFailed      = newError(2304, "更新成员失败")
	ErrGetTeamMemberListFailed = newError(2305, "获取团队成员列表失败")
	ErrAddTeamMemberFailed     = newError(2306, "添加团队成员失败")

	// 团队申请
	ErrApplyExisted         = newError(2400, "申请已存在")
	ErrApplyNotExist        = newError(2401, "申请不存在")
	ErrApplyFailed          = newError(2402, "申请失败")
	ErrGetApplyListFailed   = newError(2403, "获取申请列表失败")
	ErrApplyStatusInvalid   = newError(2404, "申请状态无效")
	ErrHandleApplyFailed    = newError(2405, "处理申请失败")
	ErrInvalidHandleStatus  = newError(2406, "无效的处理状态")
	ErrNoWithdrawPermission = newError(2407, "无撤回权限")

	// 3000 数据库
	ErrDatabase = newError(3000, "数据库错误")
	//ErrInsertFailed = newError(3001, "插入失败")
	ErrUpdateFailed = newError(3002, "更新失败")
	ErrDeleteFailed = newError(3003, "删除失败")
	ErrQueryFailed  = newError(3004, "查询失败")
	ErrDuplicateKey = newError(3005, "唯一键冲突")

	// 4000 es
	ErrCreateEsArticleFailed = newError(4000, "创建es文章失败")
	ErrUpdateEsArticleFailed = newError(4001, "更新es文章失败")
	ErrDeleteEsArticleFailed = newError(4002, "删除es文章失败")
	//ErrQueryEsArticleFailed  = newError(4003, "查询es文章失败")
	ErrCreateEsIndexFailed  = newError(4004, "创建es索引失败")
	ErrDeleteEsIndexFailed  = newError(4005, "删除es索引失败")
	ErrCreateEsMapperFailed = newError(4006, "创建es映射失败")

	// 20000 业务逻辑错误
	ErrParamEmpty          = newError(20000, "参数为空")
	ErrUserAlreadyAuth     = newError(20001, "用户已认证")
	ErrUserAuthPending     = newError(20002, "用户认证待处理")
	ErrUserAuthFailed      = newError(20003, "用户认证失败")
	ErrArticleAlreadyExist = newError(20004, "文章已存在")
	ErrCreateArticleFailed = newError(20005, "创建文章失败")
	// 团队成员
	ErrAtLeastOneAdmin = newError(20006, "至少需要一个管理员")
)
