package passport

import (
	"github.com/gin-gonic/gin"
)

func Router(app *gin.Engine, conf *Config, handlers *Handlers, middleware *IdentityMiddleware, notificationHandlers *NotificationHandlers) {
	group := app.Group("")
	{
		group.POST("/admins", middleware.Authenticate(), middleware.Authorize("admin"), handlers.CreateAdmin)
		group.POST("/users", handlers.CreateUser)
		group.GET("/users/:userId", middleware.Authenticate(), handlers.GetUserById)
		group.GET("/users", middleware.Authenticate(), middleware.Authorize("admin"), handlers.GetUsers)
		group.PATCH("/users/:userId", middleware.Authenticate(), middleware.Authorize("admin"), handlers.UpdateUser)
		group.DELETE("/users/:userId", middleware.Authenticate(), middleware.Authorize("admin"), handlers.DeleteUser)
		group.POST("/token", handlers.GetToken)
		group.POST("/verify/:token", handlers.VerifyEmail)
		group.POST("/password-recovery/email", handlers.PasswordRecovery)
		group.POST("/password-recovery/exchange", handlers.ExchangeRecoveryCode)
		group.POST("/password-recovery/reset", handlers.ResetPassword)
		group.POST("/roles", middleware.Authenticate(), middleware.Authorize("admin"), handlers.CreateRole)
		group.GET("/roles", middleware.Authenticate(), middleware.Authorize("admin"), handlers.GetRoles)
		group.PUT("/roles/:roleId", middleware.Authenticate(), middleware.Authorize("admin"), handlers.UpdateRole)
		group.POST("/rights", middleware.Authenticate(), middleware.Authorize("admin"), handlers.CreateRight)
		group.GET("/rights", middleware.Authenticate(), middleware.Authorize("admin"), handlers.GetRights)
		group.PUT("/rights/:rightId", middleware.Authenticate(), middleware.Authorize("admin"), handlers.UpdateRight)
		group.POST("/facebook/callback", handlers.FacebookCallback)
		group.GET("/notifications", notificationHandlers.Reader)
	}
}
