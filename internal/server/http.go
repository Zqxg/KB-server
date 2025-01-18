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
		// 权限包含关系：超级管理员 > 学校管理员 > 学生用户 > 普通用户
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

			// 无模块

		}
		// 学生用户路由组
		studentUserRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger, enums.SUTDENT_USER))
		{
			// 文章模块
			studentUserRouter.POST(enums.ARTICLE+"/create", articleHandler.CreateArticle)                 // 新建文章
			studentUserRouter.GET(enums.ARTICLE+"/getArticleCategory", articleHandler.GetArticleCategory) // 获取文章分组
			studentUserRouter.GET(enums.ARTICLE+"/getArticle", articleHandler.GetArticle)                 // 获取文章详细

		}
		// 学校管理员路由组
		//schoolAdminRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger, enums.SCHOOL_ADMIN))
		//{
		//}
		// 超级管理员路由组
		//superAdminRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger, enums.SUPER_ADMIN))
		//{
		//}
	}

	return s
}
