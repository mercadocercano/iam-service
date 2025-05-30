package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"iam/src/user/application/request"
	"iam/src/user/application/usecase"
	"iam/src/user/domain/exception"
	"iam/src/user/infrastructure/criteria"
)

type UserHandler struct {
	createUserUseCase          *usecase.CreateUserUseCase
	updateUserUseCase          *usecase.UpdateUserUseCase
	getUserByIDUseCase         *usecase.GetUserByIDUseCase
	listUsersUseCase           *usecase.ListUsersUseCase
	listUsersByCriteriaUseCase *usecase.ListUsersByCriteriaUseCase
	deleteUserUseCase          *usecase.DeleteUserUseCase
}

func NewUserHandler(
	createUserUseCase *usecase.CreateUserUseCase,
	updateUserUseCase *usecase.UpdateUserUseCase,
	getUserByIDUseCase *usecase.GetUserByIDUseCase,
	listUsersUseCase *usecase.ListUsersUseCase,
	listUsersByCriteriaUseCase *usecase.ListUsersByCriteriaUseCase,
	deleteUserUseCase *usecase.DeleteUserUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUseCase:          createUserUseCase,
		updateUserUseCase:          updateUserUseCase,
		getUserByIDUseCase:         getUserByIDUseCase,
		listUsersUseCase:           listUsersUseCase,
		listUsersByCriteriaUseCase: listUsersByCriteriaUseCase,
		deleteUserUseCase:          deleteUserUseCase,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param request body request.CreateUserRequest true "Create user request"
// @Success 201 {object} response.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	user, err := h.createUserUseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case exception.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "El usuario ya existe"})
		case exception.ErrInvalidEmail:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email inválido"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	user, err := h.getUserByIDUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		switch err {
		case exception.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body request.UpdateUserRequest true "Update user request"
// @Success 200 {object} response.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	req.ID = id

	user, err := h.updateUserUseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case exception.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		case exception.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "El email ya está en uso"})
		case exception.ErrInvalidEmail:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email inválido"})
		case exception.ErrInvalidStatus:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Estado inválido"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// ListUsers godoc
// @Summary List users with advanced criteria
// @Description List users with filtering, sorting and pagination using advanced criteria
// @Tags users
// @Accept json
// @Produce json
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param tenant_id query string false "Tenant ID (optional if header provided)"
// @Param email query string false "Filter by email (partial search)"
// @Param first_name query string false "Filter by first name (partial search)"
// @Param last_name query string false "Filter by last name (partial search)"
// @Param status query string false "Filter by user status"
// @Param role_id query string false "Filter by role ID"
// @Param provider query string false "Filter by authentication provider"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param sort_by query string false "Sort field" default("created_at")
// @Param sort_dir query string false "Sort direction (asc/desc)" default("desc")
// @Success 200 {object} response.UserListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Verificar que el header X-Tenant-ID esté presente
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Tenant-ID header es requerido"})
		return
	}

	// Agregar tenant_id a los query parameters para el filtrado automático
	query := c.Request.URL.Query()
	query.Set("tenant_id", tenantID)
	c.Request.URL.RawQuery = query.Encode()

	// Utilizar el criteria builder para construir los criterios desde la petición
	criteriaBuilder := criteria.NewUserCriteriaBuilder()
	crit := criteriaBuilder.BuildValidated(c)

	// Ejecutar el caso de uso para listar usuarios
	users, err := h.listUsersByCriteriaUseCase.Execute(c.Request.Context(), crit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	err = h.deleteUserUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		switch err {
		case exception.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor", "details": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// RegisterRoutes registra las rutas del módulo user
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("", h.CreateUser)
		userGroup.GET("/:id", h.GetUserByID)
		userGroup.PUT("/:id", h.UpdateUser)
		userGroup.DELETE("/:id", h.DeleteUser)
		userGroup.GET("", h.ListUsers)
	}
}
