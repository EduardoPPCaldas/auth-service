package usecases

import (
	"context"
	"fmt"

	"github.com/EduardoPPCaldas/auth-service/internal/application/user/services/token"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/google/uuid"
)

type LogoutUseCase interface {
	Execute(ctx context.Context, userID string) error
	LogoutSingle(ctx context.Context, refreshToken string) error
}

type logoutUseCase struct {
	userRepo            user.UserRepository
	refreshTokenService token.Service
}

func NewLogoutUseCase(
	userRepo user.UserRepository,
	refreshTokenService token.Service,
) LogoutUseCase {
	return &logoutUseCase{
		userRepo:            userRepo,
		refreshTokenService: refreshTokenService,
	}
}

func (uc *logoutUseCase) Execute(ctx context.Context, userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	if err := uc.refreshTokenService.RevokeAllUserTokens(ctx, userUUID); err != nil {
		return fmt.Errorf("failed to revoke all user tokens: %w", err)
	}

	return nil
}

func (uc *logoutUseCase) LogoutSingle(ctx context.Context, refreshToken string) error {
	refreshTokenEntity, err := uc.refreshTokenService.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("invalid refresh token: %w", err)
	}

	if err := uc.refreshTokenService.RevokeRefreshToken(ctx, refreshTokenEntity.ID); err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	return nil
}
