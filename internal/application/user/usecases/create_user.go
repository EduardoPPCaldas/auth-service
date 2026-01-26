package usecases

import (
	"errors"
	"fmt"

	"github.com/EduardoPPCaldas/auth-service/internal/application/user/services/token"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateUserUseCase struct {
	userRepository user.UserRepository
	tokenGenerator token.TokenGenerator
}

func NewCreateUserUseCase(userRepository user.UserRepository, tokenGenerator token.TokenGenerator) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: userRepository,
		tokenGenerator: tokenGenerator,
	}
}

func (u *CreateUserUseCase) Execute(email, password string) (string, error) {
	existingUser, err := u.userRepository.FindByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("error creating user: %w", err)
	}

	if existingUser != nil {
		return "", fmt.Errorf("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}

	user := user.New(email, lo.ToPtr(string(hashedPassword)))
	err = u.userRepository.Create(user)
	if err != nil {
		return "", fmt.Errorf("error creating user: %w", err)
	}

	return u.tokenGenerator.GenerateToken(user)
}
