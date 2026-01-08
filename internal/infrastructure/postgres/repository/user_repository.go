package repository

import (
	"os/user"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *user.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*user.User, error) {
	var user user.User
	return &user, nil
}