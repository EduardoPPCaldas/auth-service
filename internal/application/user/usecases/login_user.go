package usecases

import (
	"fmt"
	"os"
	"time"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserUseCase struct {
	userRepository user.UserRepository
}

func NewLoginUserUseCase(userRepository user.UserRepository) *LoginUserUseCase {
	return &LoginUserUseCase{userRepository: userRepository}
}

func (u *LoginUserUseCase) Execute(email, password string) (string, error) {
	user, err := u.userRepository.FindByEmail(email)
	if err != nil {
		return "", fmt.Errorf("error finding user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}
