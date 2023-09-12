package permissions

import (
	"net/http"

	"github.com/georgi-georgiev/blunder"
	"github.com/georgi-georgiev/passport/payloads"
	"github.com/georgi-georgiev/passport/responses"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type PermissionHandlers struct {
	roleService  *RoleService
	rightService *RightService
	log          *zap.Logger
	blunder      *blunder.Blunder
}

func NewPermissionHandlers(roleService *RoleService, rightService *RightService, log *zap.Logger, blunder *blunder.Blunder) *PermissionHandlers {
	return &PermissionHandlers{roleService: roleService, rightService: rightService, log: log, blunder: blunder}
}

// CreateRoleHandler godoc
// @Summary Create role
// @Description create role
// @Tags identity
// @Accept  json
// @Produce  json
// @Security OAuth2Application
// @Param data body CreateRolePayload true "data"
// @Router /roles [post]
func (h *PermissionHandlers) CreateRole(c *gin.Context) {
	var payload payloads.CreateRolePayload
	errors := h.blunder.BindJson(c.Request, &payload)
	if errors != nil {
		for _, err := range errors {
			h.blunder.GinAdd(c, err)
		}
		return
	}

	_, err := h.roleService.CreateRole(c.Request.Context(), payload)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// GetRolesHandler godoc
// @Summary Get roles
// @Description get roles
// @Tags identity
// @Accept  json
// @Produce  json
// @Security OAuth2Application
// @Success 200 {array} RoleResponse
// @Router /roles [get]
func (h *PermissionHandlers) GetRoles(c *gin.Context) {
	roles, err := h.roleService.GetRoles(c.Request.Context())
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	response := []responses.RoleResponse{}
	for _, r := range roles {
		response = append(response, responses.RoleResponse{
			Name: r.Name,
		})
	}

	c.JSON(http.StatusOK, response)
}

// UpdateRoleHandler godoc
// @Summary Update role
// @Description update role
// @Tags identity
// @Accept  json
// @Produce  json
// @Security OAuth2Application
// @Param roleID path string true "1"
// @Param data body UpdateRolePayload true "data"
// @Success 200 {object} RoleResponse
// @Router /roles/{roleId} [put]
func (h *PermissionHandlers) UpdateRole(c *gin.Context) {
	roleIdParam := c.Param("roleId")
	if roleIdParam == "" {
		c.JSON(http.StatusBadRequest, blunder.BadRequest())
		return
	}

	roleId, err := primitive.ObjectIDFromHex(roleIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, blunder.BadRequest())
		return
	}

	var payload payloads.UpdateRolePayload

	errors := h.blunder.BindJson(c.Request, &payload)
	if errors != nil {
		for _, err := range errors {
			h.blunder.GinAdd(c, err)
		}
		return
	}

	role, err := h.roleService.UpdateRole(c.Request.Context(), roleId, payload.Name)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	if role == nil {
		c.JSON(http.StatusNotFound, blunder.NotFound())
		return
	}

	c.JSON(http.StatusOK, responses.RoleResponse{Name: role.Name})
}

// CreateRightHandler godoc
// @Summary Create right
// @Description create right
// @Tags identity
// @Accept  json
// @Produce  json
// @Security OAuth2Application
// @Param data body CreateRightPayload true "data"
// @Router /rights [post]
func (h *PermissionHandlers) CreateRight(c *gin.Context) {
	var payload payloads.CreateRightPayload
	errors := h.blunder.BindJson(c.Request, &payload)
	if errors != nil {
		for _, err := range errors {
			h.blunder.GinAdd(c, err)
		}
		return
	}

	_, err := h.rightService.CreateRight(c.Request.Context(), payload)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// GetRightsHandler godoc
// @Summary Get rights
// @Description get rights
// @Tags identity
// @Accept  json
// @Produce  json
// @Security OAuth2Application
// @Success 200 {array} RightResponse
// @Router /rights [get]
func (h *PermissionHandlers) GetRights(c *gin.Context) {
	rights, err := h.rightService.GetRights(c.Request.Context())
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	response := []responses.RightResponse{}
	for _, r := range rights {
		response = append(response, responses.RightResponse{
			Name: r.Name,
		})
	}

	c.JSON(http.StatusOK, response)
}

// UpdateRoleHandler godoc
// @Summary Update role
// @Description update role
// @Tags identity
// @Accept  json
// @Produce  json
// @Security OAuth2Application
// @Param rightID path string true "1"
// @Param data body UpdateRightPayload true "data"
// @Success 200 {object} RightResponse
// @Router /rights/{rightId} [put]
func (h *PermissionHandlers) UpdateRight(c *gin.Context) {
	rightIdParam := c.Param("rightId")
	if rightIdParam == "" {
		c.JSON(http.StatusBadRequest, blunder.BadRequest())
		return
	}

	rightId, err := primitive.ObjectIDFromHex(rightIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, blunder.BadRequest())
		return
	}

	var payload payloads.UpdateRightPayload

	errors := h.blunder.BindJson(c.Request, &payload)
	if errors != nil {
		for _, err := range errors {
			h.blunder.GinAdd(c, err)
		}
		return
	}

	right, err := h.rightService.UpdateRight(c.Request.Context(), rightId, payload.Name)
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	if right == nil {
		c.JSON(http.StatusNotFound, blunder.NotFound())
		return
	}

	c.JSON(http.StatusOK, responses.RightResponse{Name: right.Name})
}
