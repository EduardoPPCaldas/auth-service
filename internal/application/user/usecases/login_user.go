package usecases

import (
	"context"
	"fmt"

	"github.com/EduardoPPCaldas/auth-service/internal/application/user/services/token"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserUseCase struct {
	userRepository user.UserRepository
	tokenGenerator token.TokenGenerator
}

func NewLoginUserUseCase(userRepository user.UserRepository, tokenGenerator token.TokenGenerator) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepository: userRepository,
		tokenGenerator: tokenGenerator,
	}
}

func (u *LoginUserUseCase) Execute(ctx context.Context, email, password string) (string, error) {
	user, err := u.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("error finding user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}

	return u.tokenGenerator.GenerateToken(user)
}
