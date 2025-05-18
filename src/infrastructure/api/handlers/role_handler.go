package handlers

import (
	"iam/src/application/services"
	"iam/src/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoleHandler struct {
	roleService *services.RoleService
}

func NewRoleHandler(roleService *services.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// CreateRole godoc
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.roleService.CreateRole(c.Request.Context(), &role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// UpdateRole godoc
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role.ID = id
	if err := h.roleService.UpdateRole(c.Request.Context(), &role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// DeleteRole godoc
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.roleService.DeleteRole(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetRoleByID godoc
func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	role, err := h.roleService.GetRoleByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// GetAllRoles godoc
func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	c.Set("log", "GetAllRoles handler called")
	roles, err := h.roleService.GetAllRoles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// GetRolesBySaas godoc
func (h *RoleHandler) GetRolesBySaas(c *gin.Context) {
	saasType := models.SaasType(c.Query("saas"))
	if !saasType.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de SaaS inválido"})
		return
	}

	roles, err := h.roleService.GetRolesBySaas(c.Request.Context(), saasType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// GetRoleByName godoc
func (h *RoleHandler) GetRoleByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El nombre es requerido"})
		return
	}

	role, err := h.roleService.GetRoleByName(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}
