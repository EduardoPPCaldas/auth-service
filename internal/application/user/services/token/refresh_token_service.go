package token

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/token"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/google/uuid"
)

type Service interface {
	GenerateRefreshToken(ctx context.Context, user *user.User) (string, error)
	ValidateRefreshToken(ctx context.Context, tokenString string) (*token.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenID uuid.UUID) error
	RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error
}

type service struct {
	tokenRepo     token.Repository
	userRepo      user.UserRepository
	refreshExpiry time.Duration
}

func NewRefreshTokenService(tokenRepo token.Repository, userRepo user.UserRepository, refreshExpiry time.Duration) Service {
	return &service{
		tokenRepo:     tokenRepo,
		userRepo:      userRepo,
		refreshExpiry: refreshExpiry,
	}
}

func (s *service) GenerateRefreshToken(ctx context.Context, user *user.User) (string, error) {
	tokenString, err := generateRandomToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	tokenHash := hashToken(tokenString)

	refreshToken := &token.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.refreshExpiry),
		CreatedAt: time.Now(),
	}

	if err := s.tokenRepo.Create(ctx, refreshToken); err != nil {
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return tokenString, nil
}

func (s *service) ValidateRefreshToken(ctx context.Context, tokenString string) (*token.RefreshToken, error) {
	tokenHash := hashToken(tokenString)

	refreshToken, err := s.tokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found: %w", err)
	}

	if !refreshToken.IsValid() {
		return nil, fmt.Errorf("refresh token is invalid or expired")
	}

	return refreshToken, nil
}

func (s *service) RevokeRefreshToken(ctx context.Context, tokenID uuid.UUID) error {
	return s.tokenRepo.Revoke(ctx, tokenID)
}

func (s *service) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	return s.tokenRepo.RevokeByUserID(ctx, userID)
}

func generateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
