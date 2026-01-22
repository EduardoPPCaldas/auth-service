package http

import (
	"github.com/EduardoPPCaldas/auth-service/internal/presentation/http/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupRoutes configures all routes for the Echo instance
func SetupRoutes(e *echo.Echo, authHandler *handlers.AuthHandler) {
	// Swagger/OpenAPI documentation
	swaggerHandler := handlers.NewSwaggerHandler()
	e.GET("/swagger.json", swaggerHandler.GetSwaggerJSON)
	e.GET("/swagger", swaggerHandler.GetSwaggerUI)
	e.GET("/swagger/", swaggerHandler.GetSwaggerUI)

	// Routes
	v1 := e.Group("/api/v1")
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
}

// SetupMiddleware configures middleware for the Echo instance
func SetupMiddleware(e *echo.Echo) {
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
}
