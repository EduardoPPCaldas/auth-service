package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/EduardoPPCaldas/auth-service/internal/application/user/services/oauth"
	"github.com/EduardoPPCaldas/auth-service/internal/application/user/services/token"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/role"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"gorm.io/gorm"
)

type LoginWithGoogleUseCase struct {
	userRepository  user.UserRepository
	roleRepository  role.Repository
	tokenGenerator  token.TokenGenerator
	googleValidator oauth.GoogleTokenValidator
}

func NewLoginWithGoogleUseCase(
	userRepository user.UserRepository,
	roleRepository role.Repository,
	tokenGenerator token.TokenGenerator,
	googleValidator oauth.GoogleTokenValidator,
) *LoginWithGoogleUseCase {
	return &LoginWithGoogleUseCase{
		userRepository:  userRepository,
		roleRepository:  roleRepository,
		tokenGenerator:  tokenGenerator,
		googleValidator: googleValidator,
	}
}

func (u *LoginWithGoogleUseCase) Execute(ctx context.Context, idToken string) (string, error) {
	googleUser, err := u.googleValidator.Validate(ctx, idToken)
	if err != nil {
		return "", fmt.Errorf("invalid google token: %w", err)
	}

	var appUser *user.User

	existingUser, err := u.userRepository.FindByEmail(googleUser.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newUser := user.New(googleUser.Email, nil)

		// Only assign role if RBAC is enabled
		if u.roleRepository.IsRBACEnabled() {
			defaultRole, err := u.roleRepository.FindOrCreateDefault()
			if err != nil {
				return "", fmt.Errorf("failed to find default role: %w", err)
			}
			if defaultRole != nil {
				newUser.RoleID = &defaultRole.ID
				newUser.Role = defaultRole
			}
		}

		if err := u.userRepository.Create(newUser); err != nil {
			return "", fmt.Errorf("failed to create user: %w", err)
		}
		appUser = newUser
	} else if err != nil {
		return "", fmt.Errorf("failed to find user: %w", err)
	} else if existingUser != nil {
		appUser = existingUser
	}

	return u.tokenGenerator.GenerateToken(appUser)
}
