package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/EduardoPPCaldas/auth-service/internal/application/user/services/token"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/role"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateUserUseCase struct {
	userRepository user.UserRepository
	roleRepository role.Repository
	tokenGenerator token.TokenGenerator
}

func NewCreateUserUseCase(userRepository user.UserRepository, roleRepository role.Repository, tokenGenerator token.TokenGenerator) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: userRepository,
		roleRepository: roleRepository,
		tokenGenerator: tokenGenerator,
	}
}

func (u *CreateUserUseCase) Execute(ctx context.Context, email, password string) (string, error) {
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

	newUser := user.New(email, lo.ToPtr(string(hashedPassword)))

	// Only assign role if RBAC is enabled
	if u.roleRepository.IsRBACEnabled(ctx) {
		defaultRole, err := u.roleRepository.FindOrCreateDefault(ctx)
		if err != nil {
			return "", fmt.Errorf("error finding default role: %w", err)
		}
		if defaultRole != nil {
			newUser.RoleID = &defaultRole.ID
			newUser.Role = defaultRole
		}
	}

	err = u.userRepository.Create(newUser)
	if err != nil {
		return "", fmt.Errorf("error creating user: %w", err)
	}

	return u.tokenGenerator.GenerateToken(newUser)
}
