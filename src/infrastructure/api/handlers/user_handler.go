package handlers

import (
	"iam/src/application/services"
	"iam/src/domain/models"
	"iam/src/domain/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Email    string            `json:"email" binding:"required,email"`
	Password string            `json:"password" binding:"required,min=8"`
	TenantID uuid.UUID         `json:"tenant_id" binding:"required"`
	RoleID   uuid.UUID         `json:"role_id" binding:"required"`
	Status   models.UserStatus `json:"status"`
}

type UpdateUserRequest struct {
	Email    string            `json:"email" binding:"omitempty,email"`
	RoleID   interface{}       `json:"role_id,omitempty"` // Aceptar string o UUID
	Status   models.UserStatus `json:"status,omitempty"`
	Provider string            `json:"-"` // Ignorar este campo
}

type UserHandler struct {
	userService *services.UserService
	roleRepo    repositories.RoleRepository
}

func NewUserHandler(userService *services.UserService, roleRepo repositories.RoleRepository) *UserHandler {
	return &UserHandler{
		userService: userService,
		roleRepo:    roleRepo,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		Email:    req.Email,
		TenantID: req.TenantID,
		RoleID:   req.RoleID,
		Status:   req.Status,
	}

	if err := h.userService.CreateUser(c.Request.Context(), user, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// No devolvemos el hash de la contraseña en la respuesta
	user.PasswordHash = ""
	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// No devolvemos el hash de la contraseña en la respuesta
	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	var tenantID *uuid.UUID
	if tenantIDStr := c.Query("tenant_id"); tenantIDStr != "" {
		parsed, err := uuid.Parse(tenantIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant_id format"})
			return
		}
		tenantID = &parsed
	}

	user, err := h.userService.GetUserByEmail(c.Request.Context(), email, tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// No devolvemos el hash de la contraseña en la respuesta
	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUsersByTenant(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant_id format"})
		return
	}

	users, err := h.userService.GetUsersByTenant(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// No devolvemos el hash de la contraseña en la respuesta
	for i := range users {
		users[i].PasswordHash = ""
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar que el usuario existe
	_, err = h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Crear un objeto UserUpdate con los campos proporcionados
	update := &models.UserUpdate{
		ID: id,
	}

	if req.Email != "" {
		update.Email = req.Email
	}
	if req.RoleID != nil {
		// Intentar obtener el rol por nombre si es string
		if roleName, ok := req.RoleID.(string); ok {
			role, err := h.roleRepo.GetByName(c.Request.Context(), roleName)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role name"})
				return
			}
			update.RoleID = role.ID
		} else if roleID, ok := req.RoleID.(uuid.UUID); ok {
			update.RoleID = roleID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role_id format"})
			return
		}
	}
	if req.Status != "" {
		update.Status = req.Status
	}

	if err := h.userService.UpdateUser(c.Request.Context(), update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Obtener el usuario actualizado
	updatedUser, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving updated user"})
		return
	}

	updatedUser.PasswordHash = ""
	c.JSON(http.StatusOK, updatedUser)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
