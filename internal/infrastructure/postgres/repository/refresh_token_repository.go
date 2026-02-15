package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/token"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefreshTokenRepository implements token.Repository interface
type RefreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository creates a new refresh token repository
func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create creates a new refresh token
func (r *RefreshTokenRepository) Create(ctx context.Context, t *token.RefreshToken) error {
	return r.db.WithContext(ctx).Create(t).Error
}

// FindByTokenHash finds a refresh token by its hash
func (r *RefreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*token.RefreshToken, error) {
	var rt token.RefreshToken
	err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

// FindByUserID finds all refresh tokens for a user
func (r *RefreshTokenRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*token.RefreshToken, error) {
	var tokens []*token.RefreshToken
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&tokens).Error
	return tokens, err
}

// Revoke revokes a specific refresh token
func (r *RefreshTokenRepository) Revoke(ctx context.Context, tokenID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&token.RefreshToken{}).
		Where("id = ?", tokenID).
		Update("revoked_at", &now).Error
}

// RevokeByUserID revokes all refresh tokens for a user
func (r *RefreshTokenRepository) RevokeByUserID(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&token.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", &now).Error
}

// CleanExpired removes expired refresh tokens
func (r *RefreshTokenRepository) CleanExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&token.RefreshToken{}).Error
}

// HashToken creates a SHA256 hash of a token
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
