package main

import (
	"log"
	"iam/src/application/services"
	"iam/src/infrastructure/api/handlers"
	"iam/src/infrastructure/api/routes"
	"iam/src/infrastructure/config"
	"iam/src/infrastructure/persistence"

	"github.com/gin-gonic/gin"
)

func main() {
	// Configuración de la base de datos
	dbConfig := config.NewDatabaseConfig()
	db, err := config.NewDatabaseConnection(dbConfig)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Inicialización de repositorios
	planRepo := persistence.NewPostgresPlanRepository(db)
	tenantRepo := persistence.NewPostgresTenantRepository(db)
	userRepo := persistence.NewPostgresUserRepository(db)
	roleRepo := persistence.NewPostgresRoleRepository(db)
	authRepo := persistence.NewPostgresAuthRepository(db)

	// Inicialización de servicios
	planService := services.NewPlanService(planRepo)
	tenantService := services.NewTenantService(tenantRepo)
	userService := services.NewUserService(userRepo)
	roleService := services.NewRoleService(roleRepo)

	// Configuración e inicialización del servicio de autenticación
	authConfig := config.NewAuthConfig()
	authService, err := services.NewAuthService(services.AuthConfig{
		JWTSecret:          authConfig.JWTSecret,
		AccessTokenExpiry:  authConfig.AccessTokenExpiry,
		RefreshTokenExpiry: authConfig.RefreshTokenExpiry,
		GoogleClientID:     authConfig.GoogleClientID,
	}, userRepo, authRepo)
	if err != nil {
		log.Fatalf("Error initializing auth service: %v", err)
	}

	// Inicialización de handlers
	planHandler := handlers.NewPlanHandler(planService)
	tenantHandler := handlers.NewTenantHandler(tenantService, userService, roleRepo)
	userHandler := handlers.NewUserHandler(userService, roleRepo)
	authHandler := handlers.NewAuthHandler(authService)
	roleHandler := handlers.NewRoleHandler(roleService)

	// Configuración del router
	router := gin.Default()

	// Configuración de CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Configuración de rutas
	routes.SetupRoutes(router, planHandler, tenantHandler, userHandler, authHandler, roleHandler)

	// Iniciar el servidor
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
