package token

import (
	"fmt"
	"os"
	"time"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type tokenGenerator struct {
}

type TokenGenerator interface {
	GenerateToken(user *user.User) (string, error)
	ExtractUserID(tokenString string) (uuid.UUID, error)
}

func NewTokenGenerator() TokenGenerator {
	return &tokenGenerator{}
}

func (t *tokenGenerator) GenerateToken(user *user.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable is not set")
	}

	claims := jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	// Only add role and permissions if user has a role (RBAC is enabled)
	if user.Role != nil {
		claims["role"] = user.Role.Name
		claims["permissions"] = user.Role.GetPermissionStrings()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}

func (t *tokenGenerator) ExtractUserID(tokenString string) (uuid.UUID, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return uuid.Nil, fmt.Errorf("JWT_SECRET environment variable is not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["sub"].(string)
		if !ok {
			return uuid.Nil, fmt.Errorf("user ID not found in token claims")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, fmt.Errorf("invalid user ID format: %w", err)
		}

		return userID, nil
	}

	return uuid.Nil, fmt.Errorf("invalid token")
}
