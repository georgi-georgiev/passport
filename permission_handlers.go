package passport

import (
	"net/http"

	"github.com/georgi-georgiev/blunder"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateRoleHandler godoc
// @Summary Create role
// @Description create role
// @Tags identity
// @Accept  json
// @Produce  json
// @Security OAuth2Application
// @Param data body CreateRolePayload true "data"
// @Router /roles [post]
func (h *Handlers) CreateRole(c *gin.Context) {
	var payload CreateRolePayload
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
func (h *Handlers) GetRoles(c *gin.Context) {
	roles, err := h.roleService.GetRoles(c.Request.Context())
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	response := []RoleResponse{}
	for _, r := range roles {
		response = append(response, RoleResponse{
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
func (h *Handlers) UpdateRole(c *gin.Context) {
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

	var payload UpdateRolePayload

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

	c.JSON(http.StatusOK, RoleResponse{Name: role.Name})
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
func (h *Handlers) CreateRight(c *gin.Context) {
	var payload CreateRightPayload
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
func (h *Handlers) GetRights(c *gin.Context) {
	rights, err := h.rightService.GetRights(c.Request.Context())
	if err != nil {
		h.blunder.GinAdd(c, err)
		return
	}

	response := []RightResponse{}
	for _, r := range rights {
		response = append(response, RightResponse{
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
func (h *Handlers) UpdateRight(c *gin.Context) {
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

	var payload UpdateRightPayload

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

	c.JSON(http.StatusOK, RightResponse{Name: right.Name})
}
