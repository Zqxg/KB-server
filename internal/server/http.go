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
			commonUserRouter.GET(enums.USER+"/logout", userHandler.Logout)                    // 退出
			commonUserRouter.GET(enums.USER+"/cancel", userHandler.Cancel)                    // 注销
			commonUserRouter.GET(enums.USER+"/getUserInfo", userHandler.GetUserInfo)          // 获取用户信息
			commonUserRouter.POST(enums.USER+"/updateProfile", userHandler.UpdateProfile)     // 修改用户信息
			commonUserRouter.GET(enums.USER+"/getCollege", collegeHandler.GetCollege)         // 获取学院信息
			commonUserRouter.GET(enums.USER+"/getCollegeList", collegeHandler.GetCollegeList) // 获取学院信息列表
			commonUserRouter.POST(enums.USER+"/userAuth", userHandler.UserAuth)

			// 文章模块
			commonUserRouter.GET(enums.ARTICLE+"/getArticleCategory", articleHandler.GetArticleCategory)             // 获取文章分组
			commonUserRouter.GET(enums.ARTICLE+"/getArticle", articleHandler.GetArticle)                             // 获取文章详细
			commonUserRouter.GET(enums.ARTICLE+"/getArticleListByCategory", articleHandler.GetArticleListByCategory) // 分类获取公开文章列表
			commonUserRouter.POST(enums.ARTICLE+"/getArticleListByEs", articleHandler.GetArticleListByEs)            // es文章查询
			commonUserRouter.POST(enums.ARTICLE+"/create", articleHandler.CreateArticle)                             // 新建文章
			commonUserRouter.POST(enums.ARTICLE+"/updateArticle", articleHandler.UpdateArticle)                      // 修改文章
			commonUserRouter.POST(enums.ARTICLE+"/deleteArticle", articleHandler.DeleteArticle)                      // 删除文章
			commonUserRouter.POST(enums.ARTICLE+"/deleteArticleList", articleHandler.DeleteArticleList)              // 批量删除文章
			commonUserRouter.POST(enums.ARTICLE+"/getUserArticleList", articleHandler.GetUserArticleList)            // 获取个人文章列表

			// 知识库模块
			//commonUserRouter.GET(enums.KNOWLEDGE_BASE+"/getKBListByTeamId", knowledgeBaseHandler.GetKBListByTeamId)                      // 团队id获取知识库列表
			//commonUserRouter.GET(enums.KNOWLEDGE_BASE+"/getKnowledgeBase", knowledgeBaseHandler.GetKnowledgeBase)                              // 获取知识库详细
			//commonUserRouter.POST(enums.KNOWLEDGE_BASE+"/createKnowledgeBase", knowledgeBaseHandler.CreateKnowledgeBase) // 新建团队知识库
			//commonUserRouter.POST(enums.KNOWLEDGE_BASE+"/updateKnowledgeBase", knowledgeBaseHandler.UpdateKnowledgeBase) // 修改团队知识库
			//commonUserRouter.POST(enums.KNOWLEDGE_BASE+"/deleteKnowledgeBase", knowledgeBaseHandler.DeleteKnowledgeBase) // 删除团队知识库
			//commonUserRouter.POST(enums.KNOWLEDGE_BASE+"/deleteKnowledgeBaseList", knowledgeBaseHandler.DeleteKnowledgeBaseList)               // 批量删除团队知识库
			//commonUserRouter.POST(enums.KNOWLEDGE_BASE+"/getUserKnowledgeBaseList", knowledgeBaseHandler.GetUserKnowledgeBase)                 // 获取私人知识库
			//commonUserRouter.POST(enums.KNOWLEDGE_BASE+"/getKnowledgeBaseListByEs", knowledgeBaseHandler.GetKnowledgeBaseListByEs)             // es团队知识库查询
			//commonUserRouter.POST(enums.KNOWLEDGE_BASE+"/getKnowledgeBaseListByCategory", knowledgeBaseHandler.GetKnowledgeBaseListByCategory) // 分类获取公开知识库列表

			// 团队模块
			commonUserRouter.POST(enums.TEAM+"/createTeam", teamHandler.CreateTeam)  // 新建团队
			commonUserRouter.POST(enums.TEAM+"/updateTeam", teamHandler.UpdateTeam)  // 修改团队
			commonUserRouter.POST(enums.TEAM+"/deleteTeam", teamHandler.DeleteTeam)  // 删除团队
			commonUserRouter.GET(enums.TEAM+"/getTeamList", teamHandler.GetTeamList) // 获取团队列表
			//commonUserRouter.GET(enums.TEAM+"/getTeam", teamHandler.GetTeam)              // todo 获取团队详细
			//commonUserRouter.POST(enums.TEAM+"/getUserTeamList", teamHandler.GetUserTeam) // todo 获取个人团队列表

			//commonUserRouter.POST(enums.TEAM+"/getTeamMemberList", teamHandler.GetTeamMemberList) // 获取团队成员列表
			//commonUserRouter.POST(enums.TEAM+"/addTeamMember", teamHandler.AddTeamMember)         // 添加团队成员
			//commonUserRouter.POST(enums.TEAM+"/deleteTeamMember", teamHandler.DeleteTeamMember)   // 删除团队成员
			//commonUserRouter.POST(enums.TEAM+"/updateTeamMember", teamHandler.UpdateTeamMember)   // 修改团队成员

		}
		//// 学生用户路由组
		//studentUserRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger, enums.SUTDENT_USER))
		//{
		//
		//}
		// 学校管理员路由组
		//schoolAdminRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger, enums.SCHOOL_ADMIN))
		//{
		//}
		// 超级管理员路由组
		//superAdminRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger, enums.SUPER_ADMIN))
		//{
		//commonUserRouter.POST(enums.TEAM+"/deleteTeamList", teamHandler.DeleteTeamList)               // 管理员批量删除团队
		//}
	}

	return s
}
