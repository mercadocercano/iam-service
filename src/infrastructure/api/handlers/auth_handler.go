package handlers

import (
	"iam/src/application/services"
	"iam/src/domain/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login maneja tanto el login local como el federado
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RefreshToken genera un nuevo access token usando un refresh token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token requerido"})
		return
	}

	resp, err := h.authService.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		switch err {
		case services.ErrInvalidToken:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		case services.ErrExpiredToken:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expirado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Logout invalida todos los refresh tokens del usuario
func (h *AuthHandler) Logout(c *gin.Context) {
	// Obtener token del header Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := h.authService.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), claims.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ValidateAuth middleware para validar tokens JWT
func (h *AuthHandler) ValidateAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("log", "ValidateAuth middleware called for path: "+c.Request.URL.Path)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := h.authService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}

		// Agregar claims al contexto para uso posterior
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("tenant_id", claims.TenantID)
		c.Set("role_id", claims.RoleID)

		c.Next()
	}
}
