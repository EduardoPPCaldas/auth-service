package http

import (
	"github.com/EduardoPPCaldas/auth-service/internal/presentation/http/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupRoutes configures all routes for the Echo instance
func SetupRoutes(
	e *echo.Echo,
	authHandler *handlers.AuthHandler,
	roleHandler *handlers.RoleHandler,
	authMiddlewareFunc func() echo.MiddlewareFunc,
	adminMiddlewareFunc func() echo.MiddlewareFunc,
) {
	// Swagger/OpenAPI documentation
	swaggerHandler := handlers.NewSwaggerHandler()
	e.GET("/swagger.json", swaggerHandler.GetSwaggerJSON)
	e.GET("/swagger", swaggerHandler.GetSwaggerUI)
	e.GET("/swagger/", swaggerHandler.GetSwaggerUI)

	// Routes
	v1 := e.Group("/api/v1")

	// Auth routes (public)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.CreateUser)
		auth.POST("/login", authHandler.LoginUser)
		auth.POST("/login/google", authHandler.LoginWithGoogle)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/logout-all", authHandler.LogoutAll)
		auth.GET("/google/challenge", authHandler.ChallengeGoogleAuth)
	}

	// Admin routes (protected)
	if authMiddlewareFunc != nil && adminMiddlewareFunc != nil {
		admin := v1.Group("/admin")
		admin.Use(authMiddlewareFunc())
		admin.Use(adminMiddlewareFunc())
		{
			// Role management
			admin.GET("/roles", roleHandler.ListRoles)
			admin.POST("/roles", roleHandler.CreateRole)
			admin.GET("/roles/:id", roleHandler.GetRole)
			admin.PUT("/roles/:id", roleHandler.UpdateRole)
			admin.DELETE("/roles/:id", roleHandler.DeleteRole)
			admin.POST("/roles/assign", roleHandler.AssignRoleToUser)
		}
	}
}

// SetupMiddleware configures middleware for the Echo instance
func SetupMiddleware(e *echo.Echo) {
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
}
