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
func (r *UserRepository) Create(u *user.User) error {
	ctx := context.Background()
	return gorm.G[user.User](r.db).Create(ctx, u)
}

// FindByEmail finds a user by their email address
func (r *UserRepository) FindByEmail(email string) (*user.User, error) {
	var u user.User
	result := r.db.Preload("Role").Preload("Role.Permissions").Where("email = ?", email).First(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	return &u, nil
}

// FindByID finds a user by their ID with role preloaded
func (r *UserRepository) FindByID(id uuid.UUID) (*user.User, error) {
	var u user.User
	result := r.db.Preload("Role").Preload("Role.Permissions").Where("id = ?", id).First(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	return &u, nil
}
