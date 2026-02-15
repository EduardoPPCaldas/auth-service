package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func ExampleUsage() {
	// Initialize the auth middleware
	authMiddleware, err := NewAuthMiddleware(
		WithJWTSecret("your-super-secret-jwt-key"),
		WithTokenValidation(TokenValidation{
			RequiredIssuer:   "my-service",
			RequiredAudience: []string{"my-api"},
		}),
	)
	if err != nil {
		log.Fatal("Failed to create auth middleware:", err)
	}

	// Example with Echo framework
	e := echo.New()

	// Protected routes
	api := e.Group("/api")
	api.Use(authMiddleware.EchoMiddleware())

	// Route requiring authentication
	api.GET("/profile", func(c echo.Context) error {
		user, ok := GetUserFromEchoContext(c)
		if !ok {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "User not found"})
		}

		return c.JSON(http.StatusOK, map[string]any{
			"user_id":     user.UserID.String(),
			"permissions": user.Permissions,
			"roles":       user.Roles,
		})
	})

	// Route requiring specific permission
	adminGroup := api.Group("/admin")
	adminGroup.Use(authMiddleware.EchoRequirePermission("admin:access"))
	adminGroup.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Admin access granted"})
	})

	// Route requiring specific role
	moderatorGroup := api.Group("/moderator")
	moderatorGroup.Use(authMiddleware.EchoRequireRole("moderator"))
	moderatorGroup.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Moderator access granted"})
	})

	// Token creation examples
	tokenCreationExamples(authMiddleware)
}

func tokenCreationExamples(authMiddleware *AuthMiddleware) {
	userID := uuid.New()

	// Create basic token
	token, err := authMiddleware.CreateTokenWithDefaults(userID)
	if err != nil {
		log.Printf("Failed to create token: %v", err)
		return
	}
	fmt.Printf("Basic token: %s\n", token)

	// Create token with custom claims
	customClaims := map[string]any{
		"permissions": []string{"read:users", "write:users"},
		"roles":       []string{"admin", "moderator"},
	}

	tokenWithClaims, err := authMiddleware.CreateToken(userID, time.Now().Add(24*time.Hour), customClaims)
	if err != nil {
		log.Printf("Failed to create token with claims: %v", err)
		return
	}
	fmt.Printf("Token with claims: %s\n", tokenWithClaims)

	// Create service token
	serviceToken, err := authMiddleware.CreateServiceToken("user-service", time.Now().Add(7*24*time.Hour))
	if err != nil {
		log.Printf("Failed to create service token: %v", err)
		return
	}
	fmt.Printf("Service token: %s\n", serviceToken)

	// Validate token examples
	validateTokenExamples(authMiddleware, token)
}

func validateTokenExamples(authMiddleware *AuthMiddleware, token string) {
	// Validate token and get claims
	claims, err := authMiddleware.ValidateTokenString(token)
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		return
	}
	fmt.Printf("Valid token for user: %s\n", claims.UserID.String())

	// Extract specific information from token
	userID, err := authMiddleware.ExtractUserID(token)
	if err != nil {
		log.Printf("Failed to extract user ID: %v", err)
		return
	}
	fmt.Printf("User ID: %s\n", userID.String())

	permissions, err := authMiddleware.ExtractPermissions(token)
	if err != nil {
		log.Printf("Failed to extract permissions: %v", err)
		return
	}
	fmt.Printf("Permissions: %v\n", permissions)

	roles, err := authMiddleware.ExtractRoles(token)
	if err != nil {
		log.Printf("Failed to extract roles: %v", err)
		return
	}
	fmt.Printf("Roles: %v\n", roles)

	// Refresh token
	newToken, err := authMiddleware.RefreshToken(token, time.Now().Add(24*time.Hour))
	if err != nil {
		log.Printf("Failed to refresh token: %v", err)
		return
	}
	fmt.Printf("Refreshed token: %s\n", newToken)
}

// Example service-to-service authentication
func ExampleServiceAuth() {
	authMiddleware, err := NewAuthMiddleware(
		WithJWTSecret("inter-service-secret"),
	)
	if err != nil {
		log.Fatal("Failed to create auth middleware:", err)
	}

	// Create token for service A
	serviceToken, err := authMiddleware.CreateServiceToken("service-a", time.Now().Add(1*time.Hour))
	if err != nil {
		log.Fatal("Failed to create service token:", err)
	}

	// Validate service token
	serviceName, err := authMiddleware.ValidateServiceToken(serviceToken)
	if err != nil {
		log.Fatal("Failed to validate service token:", err)
	}

	fmt.Printf("Valid service token for: %s\n", serviceName)
}
