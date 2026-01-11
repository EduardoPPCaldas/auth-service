package token

import (
	"fmt"
	"os"
	"time"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/golang-jwt/jwt/v5"
)

type tokenGenerator struct {
}

type TokenGenerator interface {
	GenerateToken(user *user.User) (string, error)
}

func NewTokenGenerator() TokenGenerator {
	return &tokenGenerator{}
}

func (t *tokenGenerator) GenerateToken(user *user.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}
