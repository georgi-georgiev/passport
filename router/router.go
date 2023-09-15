package router

import (
	"github.com/georgi-georgiev/passport/notifications"
	"github.com/georgi-georgiev/passport/permissions"
	"github.com/georgi-georgiev/passport/pkg/middlewares"
	"github.com/georgi-georgiev/passport/users"
	"github.com/gin-gonic/gin"
)

func Router(app *gin.Engine, userHandlers *users.UserHandlers, permissionHandlers *permissions.PermissionHandlers, middleware *middlewares.IdentityMiddleware, notificationHandlers *notifications.NotificationHandlers) {
	group := app.Group("")
	{
		group.POST("/admins", middleware.Authenticate(), middleware.Authorize("admin"), userHandlers.CreateAdmin)
		group.POST("/users", userHandlers.CreateUser)
		group.GET("/users/:userId", middleware.Authenticate(), userHandlers.GetUserById)
		group.GET("/users", middleware.Authenticate(), middleware.Authorize("admin"), userHandlers.GetUsers)
		group.PATCH("/users/:userId", middleware.Authenticate(), middleware.Authorize("admin"), userHandlers.UpdateUser)
		group.DELETE("/users/:userId", middleware.Authenticate(), middleware.Authorize("admin"), userHandlers.DeleteUser)
		group.POST("/token", userHandlers.GetToken)
		group.POST("/.well-known/jwks.json", userHandlers.JWKS)
		group.POST("/verify/:token", userHandlers.VerifyEmail)
		group.POST("/password-recovery/email", userHandlers.PasswordRecovery)
		group.POST("/password-recovery/exchange", userHandlers.ExchangeRecoveryCode)
		group.POST("/password-recovery/reset", userHandlers.ResetPassword)
		group.POST("/roles", middleware.Authenticate(), middleware.Authorize("admin"), permissionHandlers.CreateRole)
		group.GET("/roles", middleware.Authenticate(), middleware.Authorize("admin"), permissionHandlers.GetRoles)
		group.PUT("/roles/:roleId", middleware.Authenticate(), middleware.Authorize("admin"), permissionHandlers.UpdateRole)
		group.POST("/rights", middleware.Authenticate(), middleware.Authorize("admin"), permissionHandlers.CreateRight)
		group.GET("/rights", middleware.Authenticate(), middleware.Authorize("admin"), permissionHandlers.GetRights)
		group.PUT("/rights/:rightId", middleware.Authenticate(), middleware.Authorize("admin"), permissionHandlers.UpdateRight)
		group.POST("/facebook/callback", userHandlers.FacebookCallback)
		group.GET("/notifications", notificationHandlers.Reader)
	}
}
