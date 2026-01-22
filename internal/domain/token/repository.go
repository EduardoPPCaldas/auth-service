package token

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, token *RefreshToken) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*RefreshToken, error)
	Revoke(ctx context.Context, tokenID uuid.UUID) error
	RevokeByUserID(ctx context.Context, userID uuid.UUID) error
	CleanExpired(ctx context.Context) error
}
