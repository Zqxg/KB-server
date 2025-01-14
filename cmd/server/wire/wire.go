//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"projectName/internal/handler"
	"projectName/internal/job"
	"projectName/internal/repository"
	"projectName/internal/server"
	"projectName/internal/service"
	"projectName/internal/service/article"
	"projectName/internal/service/user"
	"projectName/pkg/app"
	"projectName/pkg/jwt"
	"projectName/pkg/log"
	"projectName/pkg/server/http"
	"projectName/pkg/sid"
	"time"
)

// ProvideCaptchaExpireDuration 用于提供 time.Duration 类型的实例
func ProvideCaptchaExpireDuration() time.Duration {
	return 1 * time.Minute // 假设验证码有效期为 1 分钟
}

// 提供 repository 层的实例
var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
	repository.NewCollegeRepository,
	repository.NewArticleRepository,
)

// 提供 service 层的实例
var serviceSet = wire.NewSet(
	service.NewService,
	user.NewUserService,
	ProvideCaptchaExpireDuration, // 提供 time.Duration 类型实例
	user.NewCaptchaService,       // 使用 ProvideCaptchaExpireDuration 提供的 time.Duration 类型实例
	user.NewCollegeService,
	article.NewArticleService,
)

// 提供 handler 层的实例
var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewCollegeHandler,
	handler.NewArticleHandler,
)

// 提供 job 层的实例
var jobSet = wire.NewSet(
	job.NewJob,
	job.NewUserJob,
)

// 提供 server 层的实例
var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJobServer,
)

// newApp 用于构建 App 实例
func newApp(
	httpServer *http.Server,
	jobServer *server.JobServer,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, jobServer),
		app.WithName("demo-server"),
	)
}

// NewWire 是 Wire 的生成函数，用于构建 App 实例及其依赖
func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		jobSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))
}
