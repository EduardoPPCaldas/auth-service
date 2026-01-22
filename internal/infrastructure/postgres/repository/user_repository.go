package repository

import (
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(u *user.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepository) FindByEmail(email string) (*user.User, error) {
	var u user.User
	result := r.db.Where("email = ?", email).First(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	return &u, nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*user.User, error) {
	var u user.User
	result := r.db.Where("id = ?", id).First(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	return &u, nil
}
