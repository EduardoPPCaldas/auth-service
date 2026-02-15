package auth

import (
	"net/http"
	"slices"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (am *AuthMiddleware) EchoMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization header required"})
			}

			tokenString := extractTokenFromHeader(authHeader)
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Bearer token required"})
			}

			token, err := am.parseAndValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token: " + err.Error()})
			}

			claims, ok := token.Claims.(*CustomClaims)
			if !ok || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
			}

			if err := am.validateClaims(claims); err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token validation failed: " + err.Error()})
			}

			userCtx := UserContext{
				UserID: claims.UserID,
				Claims: jwt.MapClaims{
					"sub": claims.Subject,
					"iss": claims.Issuer,
					"aud": claims.Audience,
					"exp": claims.ExpiresAt,
					"nbf": claims.NotBefore,
					"iat": claims.IssuedAt,
					"jti": claims.ID,
				},
				Permissions: claims.Permissions,
				Roles:       claims.Roles,
			}

			c.Set("user", userCtx)
			c.Set("user_id", claims.UserID.String())

			return next(c)
		}
	}
}

func (am *AuthMiddleware) EchoRequirePermission(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userCtx, ok := c.Get("user").(UserContext)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User context not found"})
			}

			if !slices.Contains(userCtx.Permissions, permission) {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Insufficient permissions"})
			}

			return next(c)
		}
	}
}

func (am *AuthMiddleware) EchoRequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userCtx, ok := c.Get("user").(UserContext)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User context not found"})
			}

			if !slices.Contains(userCtx.Roles, role) {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Insufficient role"})
			}

			return next(c)
		}
	}
}

func GetUserFromEchoContext(c echo.Context) (UserContext, bool) {
	user, ok := c.Get("user").(UserContext)
	return user, ok
}

func GetUserIDFromEchoContext(c echo.Context) (string, bool) {
	userID, ok := c.Get("user_id").(string)
	return userID, ok
}

func extractTokenFromHeader(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
