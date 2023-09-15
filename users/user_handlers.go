package users

import (
	"net/http"
	"strings"

	"github.com/georgi-georgiev/blunder"
	"github.com/georgi-georgiev/passport/payloads"
	"github.com/georgi-georgiev/passport/permissions"
	"github.com/georgi-georgiev/passport/responses"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type UserHandlers struct {
	userService  *UserService
	roleService  *permissions.RoleService
	rightService *permissions.RightService
	log          *zap.Logger
	blunder      *blunder.Blunder
}

func NewUserHandlers(userService *UserService, roleService *permissions.RoleService, rightService *permissions.RightService, log *zap.Logger, blunder *blunder.Blunder) *UserHandlers {
	return &UserHandlers{userService: userService, roleService: roleService, rightService: rightService, log: log, blunder: blunder}
}

// CreateUserHandler godoc
// @Summary Create user
// @Description create user
// @Tags identity
// @Accept  json
// @Produce  json
// @Param data body CreateUserPayload true "data"
// @Success 201 {object} CreateUserResponse
// @Failure      400  {object}  blunder.HTTPErrorResponse
// @Failure      404  {object}  blunder.HTTPErrorResponse
// @Failure      500  {object}  blunder.HTTPErrorResponse
// @Router /users [post]
func (h *UserHandlers) CreateUser(c *gin.Context) {

	if c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, blunder.BadRequest())
		return
	}

	var payload payloads.CreateUserPayload

	errors := h.blunder.BindJson(c.Request, &payload)
	if errors != nil {
		for _, err := range errors {
			h.blunder.GinAdd(c, err)
		}
		return
	}

	if payload.Role == "admin" {
		c.JSON(http.StatusForbidden, blunder.Forbidden())
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), payload.Username, payload.Email, payload.Password, payload.Role, false, payload.Rights)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	userClaims := h.userService.MapToUserClaims(user)

	token, exp, err := h.userService.IssueAccessToken(userClaims)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	refreshToken, err := h.userService.IssueRefreshToken(userClaims)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	c.JSON(http.StatusCreated, responses.CreateUserResponse{ID: user.ID.Hex(), TokenType: "Bearer", AccessToken: token, RefreshToken: refreshToken, ExpiresIn: exp})
}

// VerifyEmailHandler godoc
// @Summary Verify email
// @Description verify email
// @Tags identity
// @Accept  json
// @Produce  json
// @Param token path string true "1"
// @Router /verify/{token} [post]
func (h *UserHandlers) VerifyEmail(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, blunder.BadRequest())
		return
	}

	err := h.userService.VerifyEmail(c.Request.Context(), token)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// CreateAdminHandler godoc
// @Summary Create admin
// @Description create admin
// @Tags identity
// @Accept  json
// @Produce  json
// @Param data body CreateUserPayload true "data"
// @Success 201 {object} CreateUserResponse
// @Failure      400  {object}  blunder.HTTPErrorResponse
// @Failure      404  {object}  blunder.HTTPErrorResponse
// @Failure      500  {object}  blunder.HTTPErrorResponse
// @Router /admins [post]
func (h *UserHandlers) CreateAdmin(c *gin.Context) {

	if c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, blunder.BadRequest())
		return
	}

	var payload payloads.CreateUserPayload

	errors := h.blunder.BindJson(c.Request, &payload)
	if errors != nil {
		for _, err := range errors {
			h.blunder.GinAdd(c, err)
		}
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), payload.Username, payload.Email, payload.Password, payload.Role, true, payload.Rights)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	userClaims := h.userService.MapToUserClaims(user)

	token, exp, err := h.userService.IssueAccessToken(userClaims)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	refreshToken, err := h.userService.IssueRefreshToken(userClaims)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	c.JSON(http.StatusCreated, responses.CreateUserResponse{ID: user.ID.Hex(), TokenType: "Bearer", AccessToken: token, RefreshToken: refreshToken, ExpiresIn: exp})
}

// UpdateUserHandler godoc
// @Summary Update user
// @Description update user
// @Tags identity
// @Accept  json
// @Produce  json
// @Security OAuth2Application
// @Param userID path string true "1"
// @Param data body UpdateUserPayload true "data"
// @Success 200 {object} UserResponse
// @Router /users/{userId} [patch]
func (h *UserHandlers) UpdateUser(c *gin.Context) {
	userIDParam := c.Param("userId")
	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, blunder.BadRequest())
		return
	}

	userId, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, blunder.BadRequest())
		return
	}

	var payload payloads.UpdateUserPayload

	errors := h.blunder.BindJson(c.Request, &payload)
	if errors != nil {
		for _, err := range errors {
			h.blunder.GinAdd(c, err)
		}
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), userId, payload.Email, payload.Username, payload.IsActive, payload.ShouldChangePassword)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, blunder.NotFound())
		return
	}

	response, err := h.userService.MapToUserResponse(c.Request.Context(), user)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	c.JSON(http.StatusOK, *response)
}

// DeleteUserHandler godoc
// @Summary Delete user
// @Description delete user
// @Tags identity
// @Accept  json
// @Produce  json
// @Security OAuth2Application
// @Param userID path int true "1"
// @Router /users/{userId} [delete]
func (h *UserHandlers) DeleteUser(c *gin.Context) {
	userIDParam := c.Param("userId")
	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, blunder.BadRequest())
		return
	}

	userId, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, blunder.BadRequest())
		return
	}

	isDeleted, err := h.userService.DeleteUser(c.Request.Context(), userId)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	if !isDeleted {
		c.JSON(http.StatusNotFound, blunder.NotFound())
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// GetUsersHandler godoc
// @Summary Get users
// @Description get users
// @Tags identity
// @Accept  json
// @Produce  json
// @Security OAuth2Application
// @Success 200 {array} UserResponse
// @Router /users [get]
func (h *UserHandlers) GetUsers(c *gin.Context) {
	users, err := h.userService.GetUsers(c.Request.Context())
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	response := []responses.UserResponse{}
	for _, u := range users {
		r, err := h.userService.MapToUserResponse(c.Request.Context(), u)
		if err != nil {
			h.blunder.GinAdd(c, err)
			return
		}
		response = append(response, *r)
	}

	c.JSON(http.StatusOK, response)
}

// GetUserByIdHandler godoc
// @Summary Get user by id
// @Description get user by id
// @Tags identity
// @Accept  json
// @Produce  json
// @Security OAuth2Application
// @Param userID path int true "1"
// @Success 200 {object} UserResponse
// @Router /users/{userId} [get]
func (h *UserHandlers) GetUserById(c *gin.Context) {
	userIDParam := c.Param("userId")
	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, blunder.BadRequest())
		return
	}

	userId, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, blunder.BadRequest())
		return
	}

	user, err := h.userService.GetById(c.Request.Context(), userId)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, blunder.NotFound())
		return
	}

	response, err := h.userService.MapToUserResponse(c.Request.Context(), user)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	c.JSON(http.StatusOK, *response)
}

// GetTokenHandler godoc
// @Summary Get token
// @Description get token
// @Tags identity
// @Accept  json
// @Produce  json
// @Security BasicAuth
// @Param type query string false "refresh_token"
// @Success 200 {object} TokenResponse
// @Router /token [post]
func (h *UserHandlers) GetToken(c *gin.Context) {

	t := c.Request.URL.Query().Get("type")
	if t == "refresh_token" {
		refToken := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
		token, exp, err := h.userService.RefreshToken(c.Request.Context(), refToken)
		if err != nil {
			h.blunder.GinAdd(c, err)
		} else {
			c.JSON(http.StatusOK, responses.TokenResponse{TokenType: "Bearer", AccessToken: token, RefreshToken: refToken, ExpiresIn: exp})
		}
		return
	}

	u, p, ok := c.Request.BasicAuth()
	if ok {
		accessToken, refreshToken, exp, err := h.userService.BasicAuthToken(c.Request.Context(), u, p)
		if err != nil {
			h.blunder.GinAdd(c, err)
		} else {
			c.JSON(http.StatusOK, responses.TokenResponse{TokenType: "Bearer", AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: exp})
		}
		return
	}

	c.JSON(http.StatusBadRequest, blunder.BadRequest())
}

// GetJWKSHandler godoc
// @Summary Get jwks
// @Description get jwks
// @Tags identity
// @Accept  json
// @Produce  json
// @Security BasicAuth
// @Success 200 {object} responses.Jwks
// @Router /.well-known/jwks.json [get]
func (h *UserHandlers) JWKS(c *gin.Context) {
	key, err := h.userService.GetPublicKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, blunder.InternalServerError())
		return
	}

	jwk := responses.JSONWebKey{
		Kty: "RSA",
		Kid: "example-key",
		Use: "sig",
		N:   key.N,
		E:   key.E,
	}

	jwks := responses.Jwks{
		Keys: []responses.JSONWebKey{jwk},
	}

	c.JSON(http.StatusOK, jwks)
}

// EmailRecoveryCodeHandler godoc
// @Summary Email recovery code
// @Description email recovery code
// @Tags identity
// @Accept  json
// @Produce  json
// @Param data body RecoveryEmailPayload true "data"
// @Router /password-recovery/email [post]
func (h *UserHandlers) PasswordRecovery(c *gin.Context) {
	var payload payloads.RecoveryEmailPayload
	errors := h.blunder.BindJson(c.Request, &payload)
	if errors != nil {
		for _, err := range errors {
			h.blunder.GinAdd(c, err)
		}
		return
	}

	h.userService.SendRecoveryEmail(c.Request.Context(), payload.Email)

	c.JSON(http.StatusOK, gin.H{})
}

// ExchangeRecoveryCodeHandler godoc
// @Summary Exchange recovery code
// @Description exchange recovery code
// @Tags identity
// @Accept  json
// @Produce  json
// @Param data body ExchangeCodeRequestPayload true "data"
// @Success 200 {object} ExchangeCodeResponse
// @Router /password-recovery/exchange [post]
func (h *UserHandlers) ExchangeRecoveryCode(c *gin.Context) {
	var payload payloads.ExchangeCodeRequestPayload
	errors := h.blunder.BindJson(c.Request, &payload)
	if errors != nil {
		for _, err := range errors {
			h.blunder.GinAdd(c, err)
		}
		return
	}

	code, err := h.userService.ExchangeRecoveryCode(c.Request.Context(), payload.Email, payload.Code)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.ExchangeCodeResponse{Code: code})
}

// ResetPasswordHandler godoc
// @Summary Reset password
// @Description reset password
// @Tags identity
// @Accept  json
// @Produce  json
// @Param data body PasswordResetPayload true "data"
// @Router /password-recovery/reset [post]
func (h *UserHandlers) ResetPassword(c *gin.Context) {
	var payload payloads.PasswordResetPayload
	errors := h.blunder.BindJson(c.Request, &payload)
	if errors != nil {
		for _, err := range errors {
			h.blunder.GinAdd(c, err)
		}
		return
	}

	err := h.userService.ResetPassword(c.Request.Context(), payload.Email, payload.Code, payload.Password)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
