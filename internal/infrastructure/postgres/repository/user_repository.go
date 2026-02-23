package repository

import (
	"context"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository implements user.UserRepository interface
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user using GORM generics
func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	return gorm.G[user.User](r.db).Create(ctx, u)
}

// FindByEmail finds a user by their email address
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	user, err := gorm.G[user.User](r.db).Preload("Role.Permissions", nil).Where("email = ?", email).First(ctx)
	return &user, err
}

// FindByID finds a user by their ID with role preloaded
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	user, err := gorm.G[user.User](r.db).Preload("Role", nil).Preload("Role.Permissions", nil).Where("id = ?", id).First(ctx)
	return &user, err
}

// UpdateRole updates a user's role assignment
func (r *UserRepository) UpdateRole(ctx context.Context, userID uuid.UUID, roleID *uuid.UUID) error {
	_, err := gorm.G[user.User](r.db).Where("id = ?", userID).Update(ctx, "role_id", roleID)
	return err
}
