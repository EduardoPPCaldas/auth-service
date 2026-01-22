package middleware

import (
	"net/http"
	"strings"

	"github.com/EduardoPPCaldas/auth-service/internal/application/user/services/token"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JWTMiddleware struct {
	userRepo  user.UserRepository
	tokenGen  token.TokenGenerator
	jwtSecret string
}

type CustomClaims struct {
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}

func NewJWTMiddleware(userRepo user.UserRepository, tokenGen token.TokenGenerator, jwtSecret string) *JWTMiddleware {
	return &JWTMiddleware{
		userRepo:  userRepo,
		tokenGen:  tokenGen,
		jwtSecret: jwtSecret,
	}
}

func (m *JWTMiddleware) ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization header required"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Bearer token required"})
		}

		claims := &CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		userID, err := m.tokenGen.ExtractUserID(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
		}

		user, err := m.userRepo.FindByID(userID)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not found"})
		}

		c.Set("user", user)
		c.Set("user_id", userID)

		return next(c)
	}
}
