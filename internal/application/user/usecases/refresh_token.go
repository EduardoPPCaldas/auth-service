package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/EduardoPPCaldas/auth-service/internal/application/user/services/token"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
)

type RefreshTokenUseCase interface {
	Execute(ctx context.Context, refreshTokenString string) (*RefreshTokenResponse, error)
}

type RefreshTokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type refreshTokenUseCase struct {
	userRepo            user.UserRepository
	tokenGenerator      token.TokenGenerator
	refreshTokenService token.Service
	accessTokenExpiry   time.Duration
}

func NewRefreshTokenUseCase(
	userRepo user.UserRepository,
	tokenGenerator token.TokenGenerator,
	refreshTokenService token.Service,
	accessTokenExpiry time.Duration,
) RefreshTokenUseCase {
	return &refreshTokenUseCase{
		userRepo:            userRepo,
		tokenGenerator:      tokenGenerator,
		refreshTokenService: refreshTokenService,
		accessTokenExpiry:   accessTokenExpiry,
	}
}

func (uc *refreshTokenUseCase) Execute(ctx context.Context, refreshTokenString string) (*RefreshTokenResponse, error) {
	refreshToken, err := uc.refreshTokenService.ValidateRefreshToken(ctx, refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	user, err := uc.userRepo.FindByID(refreshToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	accessToken, err := uc.tokenGenerator.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := uc.refreshTokenService.GenerateRefreshToken(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	if err := uc.refreshTokenService.RevokeRefreshToken(ctx, refreshToken.ID); err != nil {
		return nil, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	return &RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(uc.accessTokenExpiry),
	}, nil
}
