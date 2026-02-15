# Authentication Middleware

A reusable JWT-based authentication middleware for Go services that supports multiple web frameworks.

## Features

- JWT token creation and validation
- Role-based and permission-based access control
- Support for Echo framework
- Service-to-service authentication
- Token refresh functionality
- Configurable token validation
- Custom claims support

## Installation

```bash
go get github.com/EduardoPPCaldas/auth-service/pkg/auth
```

## Quick Start

### Basic Setup

```go
package main

import (
    "log"
    "net/http"
    "time"
    
    "github.com/EduardoPPCaldas/auth-service/pkg/auth"
    "github.com/labstack/echo/v4"
    "github.com/google/uuid"
)

func main() {
    // Initialize auth middleware
    authMiddleware, err := auth.NewAuthMiddleware(
        auth.WithJWTSecret("your-super-secret-jwt-key"),
    )
    if err != nil {
        log.Fatal("Failed to create auth middleware:", err)
    }

    // Setup Echo
    e := echo.New()

    // Protected routes
    api := e.Group("/api")
    api.Use(authMiddleware.EchoMiddleware())

    api.GET("/profile", func(c echo.Context) error {
        user, ok := auth.GetUserFromEchoContext(c)
        if !ok {
            return c.JSON(http.StatusInternalServerError, map[string]string{"error": "User not found"})
        }
        return c.JSON(http.StatusOK, map[string]any{
            "user_id": user.UserID.String(),
        })
    })

    e.Logger.Fatal(e.Start(":8080"))
}
```

### Creating Tokens

```go
userID := uuid.New()

// Basic token (24h expiry)
token, err := authMiddleware.CreateTokenWithDefaults(userID)

// Token with custom claims
customClaims := map[string]any{
    "permissions": []string{"read:users", "write:users"},
    "roles":       []string{"admin"},
}

token, err := authMiddleware.CreateToken(
    userID, 
    time.Now().Add(24*time.Hour), 
    customClaims,
)

// Service token
serviceToken, err := authMiddleware.CreateServiceToken(
    "user-service", 
    time.Now().Add(7*24*time.Hour),
)
```

### Permission and Role-Based Access Control

```go
// Require specific permission
adminGroup := api.Group("/admin")
adminGroup.Use(authMiddleware.EchoRequirePermission("admin:access"))
adminGroup.GET("", func(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]string{"message": "Admin access granted"})
})

// Require specific role
moderatorGroup := api.Group("/moderator")
moderatorGroup.Use(authMiddleware.EchoRequireRole("moderator"))
moderatorGroup.GET("", func(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]string{"message": "Moderator access granted"})
})
```

### Token Validation

```go
// Validate token and get claims
claims, err := authMiddleware.ValidateTokenString(tokenString)
if err != nil {
    // Handle error (token expired, invalid, etc.)
}

// Extract specific information
userID, err := authMiddleware.ExtractUserID(tokenString)
permissions, err := authMiddleware.ExtractPermissions(tokenString)
roles, err := authMiddleware.ExtractRoles(tokenString)

// Refresh token
newToken, err := authMiddleware.RefreshToken(tokenString, time.Now().Add(24*time.Hour))
```

### Service-to-Service Authentication

```go
// Create service token
serviceToken, err := authMiddleware.CreateServiceToken("service-a", time.Now().Add(1*time.Hour))

// Validate service token
serviceName, err := authMiddleware.ValidateServiceToken(serviceToken)
if err != nil {
    // Handle invalid service token
}
```

## Configuration Options

```go
authMiddleware, err := auth.NewAuthMiddleware(
    auth.WithJWTSecret("your-secret-key"),
    auth.WithTokenValidation(auth.TokenValidation{
        SkipExpirationCheck: false,        // Set to true for testing
        RequiredIssuer:      "my-service",  // Validate issuer claim
        RequiredAudience:    []string{"my-api"}, // Validate audience claim
    }),
)
```

## User Context

The middleware injects a `UserContext` into the request context:

```go
type UserContext struct {
    UserID      uuid.UUID
    Claims      jwt.MapClaims
    Permissions []string
    Roles       []string
}
```

### Accessing User Context

```go
// In Echo handlers
user, ok := auth.GetUserFromEchoContext(c)
if ok {
    fmt.Printf("User ID: %s\n", user.UserID.String())
    fmt.Printf("Permissions: %v\n", user.Permissions)
    fmt.Printf("Roles: %v\n", user.Roles)
}

// Get just the user ID
userID, ok := auth.GetUserIDFromEchoContext(c)
```

## Error Handling

The middleware provides structured error types:

```go
var (
    auth.ErrTokenExpired     = errors.New("token has expired")
    auth.ErrTokenInvalid     = errors.New("token is invalid")
    auth.ErrTokenMissing     = errors.New("authorization token is missing")
    auth.ErrPermissionDenied = errors.New("permission denied")
    auth.ErrRoleDenied       = errors.New("role denied")
)

// Custom auth errors
authErr := auth.NewAuthError(auth.ErrorTypeExpired, "Token has expired")
// authErr.Type, authErr.Message, authErr.Code available
```

## Security Best Practices

1. **Use environment variables for secrets**:
   ```go
   secret := os.Getenv("JWT_SECRET")
   ```

2. **Set appropriate token expiry times**:
   ```go
   // Short expiry for user tokens (e.g., 1 hour)
   userTokenExpiry := time.Now().Add(1 * time.Hour)
   
   // Longer expiry for service tokens (e.g., 24 hours)
   serviceTokenExpiry := time.Now().Add(24 * time.Hour)
   ```

3. **Use HTTPS in production** to prevent token interception

4. **Validate token claims** like issuer and audience when possible

5. **Implement token revocation** for security-critical applications

## Token Structure

### User Token

```json
{
  "sub": "user-uuid",
  "permissions": ["read:users", "write:users"],
  "roles": ["admin"],
  "iss": "my-service",
  "aud": ["my-api"],
  "exp": 1640995200,
  "iat": 1640908800,
  "nbf": 1640908800
}
```

### Service Token

```json
{
  "sub": "service-name",
  "type": "service",
  "exp": 1640995200,
  "iat": 1640908800,
  "nbf": 1640908800
}
```

## Framework Support

### Echo Framework

```go
api.Use(authMiddleware.EchoMiddleware())
api.Use(authMiddleware.EchoRequirePermission("admin:access"))
api.Use(authMiddleware.EchoRequireRole("admin"))
```

### Standard Library

```go
http.Handle("/api/", authMiddleware.ValidateToken(myHandler))
http.Handle("/api/admin", authMiddleware.RequirePermission("admin:access")(myHandler))
```

## Testing

For testing, you can create tokens with known values:

```go
testUserID := uuid.New()
testToken, _ := authMiddleware.CreateTokenWithDefaults(testUserID)

// Or skip expiration checks
authMiddleware, _ := auth.NewAuthMiddleware(
    auth.WithJWTSecret("test-secret"),
    auth.WithTokenValidation(auth.TokenValidation{
        SkipExpirationCheck: true,
    }),
)
```

## Dependencies

- `github.com/golang-jwt/jwt/v5` - JWT library
- `github.com/google/uuid` - UUID generation
- `github.com/labstack/echo/v4` - Echo framework (optional)