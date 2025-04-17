package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	apiV1 "projectName/api/v1"
	"projectName/docs"
	"projectName/internal/enums"
	"projectName/internal/handler"
	"projectName/internal/middleware"
	"projectName/pkg/jwt"
	"projectName/pkg/log"
	"projectName/pkg/server/http"
)

func NewHTTPServer(
	logger *log.Logger,
	conf *viper.Viper,
	jwt *jwt.JWT,
	userHandler *handler.UserHandler,
	collegeHandler *handler.CollegeHandler,
	articleHandler *handler.ArticleHandler,
	teamHandler *handler.TeamHandler,
	adminHandler *handler.AdminHandler,
	knowledgeBaseHandler *handler.KnowledgeBaseHandler,

) *http.Server {
	gin.SetMode(gin.DebugMode)
	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)

	// swagger doc
	docs.SwaggerInfo.BasePath = "/v1"
	s.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		//ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", conf.GetInt("app.http.port"))),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
	))

	s.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
		//middleware.SignMiddleware(log),
	)
	s.GET("/", func(ctx *gin.Context) {
		logger.WithContext(ctx).Info("hello")
		apiV1.HandleSuccess(ctx, map[string]interface{}{
			":)": "Thank you for using nunu!",
		})
	})

	v1 := s.Group("/v1")
	{
		// 无需权限路由组
		noAuthRouter := v1.Group("/")
		{
			noAuthRouter.POST("/register", userHandler.Register)
			noAuthRouter.POST("/passwordLogin", userHandler.PasswordLogin)
			noAuthRouter.GET("/getCaptcha", userHandler.GetCaptcha)
		}
		// 权限包含关系：超级管理员 > 普通用户
		// 普通用户路由组
		commonUserRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger, enums.COMMON_USER))
		{
			// 用户模块
			commonUserRouter.GET(enums.USER+"/logout", userHandler.Logout)                // 退出
			commonUserRouter.GET(enums.USER+"/cancel", userHandler.Cancel)                // 注销
			commonUserRouter.GET(enums.USER+"/getUserInfo", userHandler.GetUserInfo)      // 获取用户信息
			commonUserRouter.POST(enums.USER+"/updateProfile", userHandler.UpdateProfile) // 修改用户信息
			commonUserRouter.POST(enums.USER+"/search", userHandler.Search)
			commonUserRouter.GET(enums.USER+"/getCollege", collegeHandler.GetCollege)         // 获取学院信息
			commonUserRouter.GET(enums.USER+"/getCollegeList", collegeHandler.GetCollegeList) // 获取学院信息列表
			//commonUserRouter.POST(enums.USER+"/userAuth", userHandler.UserAuth)

			// 文章模块
			commonUserRouter.GET(enums.ARTICLE+"/getArticle", articleHandler.GetArticle)                    // 获取文章详细
			commonUserRouter.POST(enums.ARTICLE+"/getArticleListByCID", articleHandler.GetArticleListByCID) // 分类获取公开文章列表
			commonUserRouter.POST(enums.ARTICLE+"/getArticleListByEs", articleHandler.GetArticleListByEs)   // es文章查询
			commonUserRouter.POST(enums.ARTICLE+"/createArticle", articleHandler.CreateArticle)             // 新建文章
			commonUserRouter.POST(enums.ARTICLE+"/updateArticle", articleHandler.UpdateArticle)             // 修改文章
			commonUserRouter.POST(enums.ARTICLE+"/deleteArticle", articleHandler.DeleteArticle)             // 删除文章
			commonUserRouter.POST(enums.ARTICLE+"/deleteArticleList", articleHandler.DeleteArticleList)     // 批量删除文章
			commonUserRouter.POST(enums.ARTICLE+"/getUserArticleList", articleHandler.GetUserArticleList)   // 获取个人文章列表

			// 团队知识库模块
			commonUserRouter.POST(enums.KNOWLEDGE_BASE+"/getKBListByTeamId", knowledgeBaseHandler.GetKBListByTeamId) // 团队id获取知识库列表
			commonUserRouter.POST(enums.KNOWLEDGE_BASE+"/getKBListByType", knowledgeBaseHandler.GetKBListByType)     // 获取知识库列表(私人知识库、公共知识库)
			commonUserRouter.GET(enums.KNOWLEDGE_BASE+"/getKBInfo", knowledgeBaseHandler.GetKBInfo)                  // 获取知识库详细
			commonUserRouter.POST(enums.KNOWLEDGE_BASE+"/createKB", knowledgeBaseHandler.CreateKB)                   // 新建团队知识库
			commonUserRouter.POST(enums.KNOWLEDGE_BASE+"/updateKBName", knowledgeBaseHandler.UpdateKBName)           // 修改团队知识库名称
			commonUserRouter.POST(enums.KNOWLEDGE_BASE+"/deleteKB", knowledgeBaseHandler.DeleteKB)                   // 删除团队知识库
			// 团队知识库分类模块
			commonUserRouter.GET(enums.KNOWLEDGE_BASE+"/getCategoryListByKB", knowledgeBaseHandler.GetCategoryListByKB) // 获取知识库分类列表
			commonUserRouter.POST(enums.KNOWLEDGE_BASE+"/createCategory", knowledgeBaseHandler.CreateCategory)          // 新建知识库分类（私人/团队）
			commonUserRouter.POST(enums.KNOWLEDGE_BASE+"/updateCategory", knowledgeBaseHandler.UpdateCategory)          // 修改知识库分类（私人/团队）
			commonUserRouter.POST(enums.KNOWLEDGE_BASE+"/deleteCategory", knowledgeBaseHandler.DeleteCategory)          // 删除知识库分类（私人/团队）

			// 团队模块
			commonUserRouter.POST(enums.TEAM+"/createTeam", teamHandler.CreateTeam)                      // 新建团队
			commonUserRouter.POST(enums.TEAM+"/updateTeam", teamHandler.UpdateTeam)                      // 修改团队
			commonUserRouter.POST(enums.TEAM+"/deleteTeam", teamHandler.DeleteTeam)                      // 删除团队
			commonUserRouter.POST(enums.TEAM+"/getTeamList", teamHandler.GetTeamList)                    // 获取团队列表
			commonUserRouter.GET(enums.TEAM+"/getTeamInfo", teamHandler.GetTeamInfo)                     // 获取团队详细
			commonUserRouter.GET(enums.TEAM+"/getUserTeamList", teamHandler.GetUserTeamList)             // 获取个人团队列表
			commonUserRouter.GET(enums.TEAM+"/getUserTeamManageList", teamHandler.GetUserTeamManageList) // 获取个人团队管理列表

			// 团队成员模块
			commonUserRouter.GET(enums.TEAM+"/getUserTeamMemberList", teamHandler.GetUserTeamMemberList) // 获取个人团队成员列表
			commonUserRouter.GET(enums.TEAM+"/getTeamMemberList", teamHandler.GetTeamMemberList)         // 获取团队成员列表
			commonUserRouter.POST(enums.TEAM+"/addTeamMember", teamHandler.AddTeamMember)                // 添加团队成员
			commonUserRouter.POST(enums.TEAM+"/deleteTeamMember", teamHandler.DeleteTeamMember)          // 删除团队成员
			commonUserRouter.POST(enums.TEAM+"/updateTeamMemberRole", teamHandler.UpdateTeamMemberRole)  // 修改团队成员角色
			commonUserRouter.POST(enums.TEAM+"/quitTeam", teamHandler.QuitTeam)                          // 退出团队
			commonUserRouter.POST(enums.TEAM+"/applyJoinTeam", teamHandler.ApplyJoinTeam)                // 申请加入团队
			commonUserRouter.POST(enums.TEAM+"/getTeamApplyList", teamHandler.GetTeamApplyList)          // 获取团队申请列表
			commonUserRouter.POST(enums.TEAM+"/handleTeamApply", teamHandler.HandleTeamApply)            // 处理团队申请

		}
		//超级管理员路由组
		superAdminRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger, enums.SUPER_ADMIN))
		{
			superAdminRouter.POST(enums.ROUTET_ADMIN+"/createPublicKB", adminHandler.CreatePublicKB)             // 新增公共知识库
			superAdminRouter.POST(enums.ROUTET_ADMIN+"/updatePublicKB", adminHandler.UpdatePublicKB)             // 更新公共知识库
			superAdminRouter.POST(enums.ROUTET_ADMIN+"/deletePublicKB", adminHandler.DeletePublicKB)             // 删除公共知识库
			superAdminRouter.POST(enums.ROUTET_ADMIN+"/createPublicCategory", adminHandler.CreatePublicCategory) // 新增公共知识库分类
			superAdminRouter.POST(enums.ROUTET_ADMIN+"/updatePublicCategory", adminHandler.UpdatePublicCategory) // 更新公共知识库分类
			superAdminRouter.POST(enums.ROUTET_ADMIN+"/deletePublicCategory", adminHandler.DeletePublicCategory) // 删除公共知识库分类

		}
	}

	return s
}
