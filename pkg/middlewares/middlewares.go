package middlewares

import (
	"net/http"
	"strings"

	"github.com/georgi-georgiev/blunder"
	"github.com/georgi-georgiev/passport/permissions"
	"github.com/georgi-georgiev/passport/users"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

type IdentityMiddleware struct {
	userRervice  *users.UserService
	roleService  *permissions.RoleService
	rightService *permissions.RightService
	log          *zap.Logger
}

func NewMiddleware(userRervice *users.UserService, roleServce *permissions.RoleService, rightService *permissions.RightService, log *zap.Logger) *IdentityMiddleware {
	return &IdentityMiddleware{userRervice: userRervice, roleService: roleServce, rightService: rightService, log: log}
}

func (m *IdentityMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		headerItems := strings.Split(authHeader, " ")

		if len(headerItems) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, blunder.Unauthorized())
			return
		}

		token := headerItems[1]
		userClaims, err := m.userRervice.ValidateToken(c.Request.Context(), token)
		if err != nil {
			m.log.With(zap.Error(err)).Error("could not validate token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, blunder.Unauthorized())
			return
		}

		userID, err := primitive.ObjectIDFromHex(userClaims.Id)
		if err != nil {
			m.log.With(zap.Error(err)).Error("could not convert user id hex to primitive")
			c.AbortWithStatusJSON(http.StatusUnauthorized, blunder.Unauthorized())
			return
		}

		user, err := m.userRervice.GetById(c.Request.Context(), userID)
		if err != nil {
			m.log.With(zap.Error(err)).Error("could not convert user id hex to primitive")
			c.AbortWithStatusJSON(http.StatusUnauthorized, blunder.Unauthorized())
			return
		}

		if user == nil {
			m.log.With(zap.Error(err)).Error("user does not exist")
			c.AbortWithStatusJSON(http.StatusUnauthorized, blunder.Unauthorized())
			return
		}

		if !user.IsActive {
			m.log.With(zap.Error(err)).Error("user is not active")
			c.AbortWithStatusJSON(http.StatusUnauthorized, blunder.Unauthorized())
			return
		}

		if !user.IsVerified {
			m.log.With(zap.Error(err)).Error("user is not verified")
			c.AbortWithStatusJSON(http.StatusUnauthorized, blunder.Unauthorized())
			return
		}

		c.Set("userId", userClaims.Id)
		c.Set("role", userClaims.Role)
		c.Set("rights", userClaims.Rights)

		c.Next()
	}
}

func (m *IdentityMiddleware) Authorize(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.GetString("role")
		roleId, err := primitive.ObjectIDFromHex(r)
		if err != nil {
			m.log.With(zap.Error(err)).Error("could not convert role id from hex to primitive object")
			c.AbortWithStatusJSON(http.StatusForbidden, blunder.Forbidden())
			return
		}

		rr, err := m.roleService.GetById(c.Request.Context(), roleId)
		if err != nil {
			m.log.With(zap.Error(err)).Error("could not get role by id")
			c.AbortWithStatusJSON(http.StatusForbidden, blunder.Forbidden())
			return
		}

		if rr == nil {
			m.log.With(zap.Error(err)).Error("role is missing")
			c.AbortWithStatusJSON(http.StatusForbidden, blunder.Forbidden())
			return
		}

		roleExists := false
		for _, role := range roles {
			if rr.Name == role {
				roleExists = true
				break
			}
		}

		if !roleExists {
			m.log.With(zap.Error(err)).Error("role is not containing in the list", zap.String("name", rr.Name), zap.Strings("roles", roles))
			c.AbortWithStatusJSON(http.StatusForbidden, blunder.Forbidden())
			return
		}

		c.Next()
	}
}

func (m *IdentityMiddleware) AuthorizeSpecific(rrs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rights := c.GetStringSlice("rights")

		rightIds := make([]primitive.ObjectID, 0)
		for _, r := range rights {
			rightId, err := primitive.ObjectIDFromHex(r)
			if err != nil {
				m.log.With(zap.Error(err)).Error("could not convert right id from hex to primitive object")
				c.AbortWithStatusJSON(http.StatusForbidden, blunder.Forbidden())
				return
			}

			rightIds = append(rightIds, rightId)
		}

		rr, err := m.rightService.GetManyByIds(c.Request.Context(), rightIds)
		if err != nil {
			m.log.With(zap.Error(err)).Error("could not get rights")
			c.AbortWithStatusJSON(http.StatusForbidden, blunder.Forbidden())
			return
		}

		rightExists := false

		for _, r := range rr {
			if slices.Contains(rrs, r.Name) {
				rightExists = true
				break
			}
		}

		if !rightExists {
			c.AbortWithStatusJSON(http.StatusForbidden, blunder.Forbidden())
			return
		}

		c.Next()
	}
}
