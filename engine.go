package passport

import (
	"github.com/georgi-georgiev/passport/docs"

	"github.com/georgi-georgiev/blunder"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func NewGinEngine(conf *Config, logger *zap.Logger, blunder *blunder.Blunder, sentryClient *sentry.Client) *gin.Engine {
	engine := gin.New()

	engine.Use(gin.CustomRecovery(blunder.GinRecovery))
	engine.Use(blunder.GinErrorHandler(logger))
	engine.HandleMethodNotAllowed = true
	engine.NoMethod(blunder.GinNoMethod)
	engine.NoRoute(blunder.GinNoRoute)

	engine.LoadHTMLFiles("blunder.html")
	engine.GET("/errors", blunder.Html)

	engine.Use(sentrygin.New(sentrygin.Options{}))

	docs.SwaggerInfo.Host = conf.Swagger.Host

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	engine.GET("/health", gin.WrapF(HandleHealthCheckWithCallback(func() string { return "Version: " + conf.App.Version })))
	engine.GET("/rstats", gin.WrapF(HandleRuntimeStats))
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	return engine
}
