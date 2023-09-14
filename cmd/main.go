package main

import (
	"context"
	"net/http"

	"github.com/georgi-georgiev/blunder"
	"github.com/georgi-georgiev/passport"
	"github.com/georgi-georgiev/passport/facade"
	"github.com/georgi-georgiev/passport/notifications"
	"github.com/georgi-georgiev/passport/permissions"
	"github.com/georgi-georgiev/passport/pkg/middlewares"
	"github.com/georgi-georgiev/passport/router"
	"github.com/georgi-georgiev/passport/users"

	_ "go.uber.org/automaxprocs"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// @title Passport
// @version 1.0
// @description Passport
// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization
// @securitydefinitions.oauth2.password OAuth2Application
// @tokenUrl http://localhost:2525//token
// @authorizationurl http://localhost:2525//token
func main() {
	fx.New(
		fx.Provide(
			passport.NewHTTPServer,
			context.Background,
			passport.NewConfig,
			passport.NewSentry,
			passport.NewLogger,
			blunder.NewRFC,
			passport.NewGinEngine,
			passport.NewMongoClient,
			passport.NewMailCleint,

			notifications.NewNotificationRepository,
			facade.NewNotificationFacade,
			notifications.NewNotificationService,
			notifications.NewNotificationHandlers,

			permissions.NewRightRepository,
			permissions.NewRightService,
			permissions.NewRoleRepository,
			permissions.NewRoleService,
			permissions.NewPermissionHandlers,

			users.NewUserRepository,
			users.NewUserService,
			users.NewUserHandlers,

			middlewares.NewMiddleware,
		),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.Invoke(func(*http.Server) {}),
		fx.Invoke(func(notifcationService *notifications.NotificationService) {
			notifcationService.Listener()
		}),
		fx.Invoke(router.Router),
	).Run()
}
