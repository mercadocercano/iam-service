package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "iam/src/application/services"
    "iam/src/domain/models"
    "iam/src/domain/repositories"
    "net/http"
)

type CreateTenantRequest struct {
    Tenant   models.Tenant `json:"tenant"`
    Password string        `json:"password" binding:"required"`
}

type TenantHandler struct {
    tenantService *services.TenantService
    userService   *services.UserService
    roleRepo      repositories.RoleRepository
}

func NewTenantHandler(tenantService *services.TenantService, userService *services.UserService, roleRepo repositories.RoleRepository) *TenantHandler {
    return &TenantHandler{
        tenantService: tenantService,
        userService:   userService,
        roleRepo:      roleRepo,
    }
}

func (h *TenantHandler) CreateTenant(c *gin.Context) {
    var req CreateTenantRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx := c.Request.Context()

    // Crear el tenant
    if err := h.tenantService.CreateTenant(ctx, &req.Tenant); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Obtener el rol ADMIN
    adminRole, err := h.roleRepo.GetByName(ctx, "ADMIN")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting admin role"})
        return
    }

    // Crear el usuario admin
    user := &models.User{
        Email:    req.Tenant.EmailUserKey,
        TenantID: req.Tenant.ID,
        RoleID:   adminRole.ID,
        Status:   models.StatusActive,
    }

    if err := h.userService.CreateUser(ctx, user, req.Password); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating admin user"})
        return
    }

    c.JSON(http.StatusCreated, req.Tenant)
}

func (h *TenantHandler) GetTenantByID(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    tenant, err := h.tenantService.GetTenantByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
        return
    }

    c.JSON(http.StatusOK, tenant)
}

func (h *TenantHandler) GetAllTenants(c *gin.Context) {
    saasType := c.Query("saas")
    var tenants []models.Tenant
    var err error

    if saasType != "" {
        tenants, err = h.tenantService.GetTenantsBySaas(c.Request.Context(), models.SaasType(saasType))
    } else {
        tenants, err = h.tenantService.GetAllTenants(c.Request.Context())
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, tenants)
}

func (h *TenantHandler) GetTenantByEmailUserKey(c *gin.Context) {
    emailUserKey := c.Query("email_user_key")
    if emailUserKey == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "email_user_key is required"})
        return
    }

    tenant, err := h.tenantService.GetTenantByEmailUserKey(c.Request.Context(), emailUserKey)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
        return
    }

    c.JSON(http.StatusOK, tenant)
}

func (h *TenantHandler) UpdateTenant(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    var tenant models.Tenant
    if err := c.ShouldBindJSON(&tenant); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    tenant.ID = id
    if err := h.tenantService.UpdateTenant(c.Request.Context(), &tenant); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, tenant)
}

func (h *TenantHandler) DeleteTenant(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    if err := h.tenantService.DeleteTenant(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}
