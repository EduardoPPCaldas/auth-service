package usecases

import (
	"errors"
	"fmt"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateUserUseCase struct {
	userRepository user.UserRepository
}

func NewCreateUserUseCase(userRepository user.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepository: userRepository}
}

func (u *CreateUserUseCase) Execute(email, password string) error {
	existingUser, err := u.userRepository.FindByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error creating user: %w", err)
	}

	if existingUser != nil {
		return fmt.Errorf("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	user := user.New(email, string(hashedPassword))
	return u.userRepository.Create(user)
}
