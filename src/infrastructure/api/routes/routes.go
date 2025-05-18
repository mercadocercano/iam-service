package routes

import (
	"iam/src/infrastructure/api/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	planHandler *handlers.PlanHandler,
	tenantHandler *handlers.TenantHandler,
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	roleHandler *handlers.RoleHandler,
) {
	// Health check endpoint (público para verificación de servicios)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "up",
			"service": "iam",
		})
	})

	// API v1 group
	v1 := router.Group("/iam/api/v1")

	// Auth routes (públicas)
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout) // Aunque requiere token, usa la ruta pública
	}

	// Middleware de autenticación para rutas protegidas
	protected := v1.Group("")
	protected.Use(authHandler.ValidateAuth())

	// Plans routes (GetAllPlans pública, resto protegidas)
	plans := v1.Group("/plans")
	{
		// Ruta pública
		plans.GET("", planHandler.GetAllPlans)

		// Rutas protegidas
		plansProtected := plans.Group("")
		plansProtected.Use(authHandler.ValidateAuth())
		{
			plansProtected.POST("", planHandler.CreatePlan)
			plansProtected.GET("/:id", planHandler.GetPlanByID)
			plansProtected.PUT("/:id", planHandler.UpdatePlan)
			plansProtected.DELETE("/:id", planHandler.DeletePlan)
		}
	}

	// Tenants routes (algunas públicas, otras protegidas)
	tenants := v1.Group("/tenants")
	{
		// Rutas públicas
		tenants.POST("", tenantHandler.CreateTenant)
		tenants.GET("/by-email-key", tenantHandler.GetTenantByEmailUserKey)

		// Rutas protegidas
		tenantsProtected := tenants.Group("")
		tenantsProtected.Use(authHandler.ValidateAuth())
		{
			tenantsProtected.GET("", tenantHandler.GetAllTenants)
			tenantsProtected.GET("/:id", tenantHandler.GetTenantByID)
			tenantsProtected.PUT("/:id", tenantHandler.UpdateTenant)
			tenantsProtected.DELETE("/:id", tenantHandler.DeleteTenant)
		}
	}

	// Users routes (protegidas)
	users := protected.Group("/users")
	{
		users.POST("", userHandler.CreateUser)
		users.GET("/by-email", userHandler.GetUserByEmail)
		users.GET("/by-tenant/:tenant_id", userHandler.GetUsersByTenant)
		users.GET("/:id", userHandler.GetUserByID)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}

	// Roles routes (todas protegidas)
	roles := protected.Group("/roles")
	{
		roles.GET("", roleHandler.GetAllRoles)
		roles.POST("", roleHandler.CreateRole)
		roles.GET("/by-saas", roleHandler.GetRolesBySaas)
		roles.GET("/by-name", roleHandler.GetRoleByName)
		roles.GET("/:id", roleHandler.GetRoleByID)
		roles.PUT("/:id", roleHandler.UpdateRole)
		roles.DELETE("/:id", roleHandler.DeleteRole)
	}
}
