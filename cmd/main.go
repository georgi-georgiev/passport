package main

import (
	"context"
	"net/http"

	"github.com/georgi-georgiev/blunder"
	"github.com/georgi-georgiev/passport"

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

			passport.NewNotificationRepository,
			passport.NewMailCleint,
			passport.NewNotificationFacade,
			passport.NewNotificationService,
			passport.NewNotificationHandlers,

			passport.NewRoleRepository,
			passport.NewUserRepository,
			passport.NewRightRepository,
			passport.NewUserService,
			passport.NewRoleService,
			passport.NewRightService,
			passport.NewHandlers,
			passport.NewMiddleware,
		),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.Invoke(func(*http.Server) {}),
		fx.Invoke(func(notifcationService *passport.NotificationService) {
			notifcationService.Listener()
		}),
		fx.Invoke(passport.Router),
	).Run()
}
