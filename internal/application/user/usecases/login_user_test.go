package usecases

import (
	"context"
	"errors"
	"testing"

	tokenmocks "github.com/EduardoPPCaldas/auth-service/internal/application/user/services/token/mocks"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginUserUseCase_Execute_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)

	useCase := NewLoginUserUseCase(mockRepo, mockTokenGen)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedPasswordStr := string(hashedPassword)
	expectedToken := "jwt-token-here"

	existingUser := &user.User{
		ID:       uuid.New(),
		Email:    email,
		Password: &hashedPasswordStr,
	}

	mockRepo.On("FindByEmail", ctx, email).Return(existingUser, nil)
	mockTokenGen.On("GenerateToken", existingUser).Return(expectedToken, nil)

	// Act
	token, err := useCase.Execute(ctx, email, password)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
	mockRepo.AssertExpectations(t)
	mockTokenGen.AssertExpectations(t)
}

func TestLoginUserUseCase_Execute_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)

	useCase := NewLoginUserUseCase(mockRepo, mockTokenGen)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	notFoundError := errors.New("user not found")

	mockRepo.On("FindByEmail", ctx, email).Return(nil, notFoundError)

	// Act
	token, err := useCase.Execute(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error finding user")
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
	mockTokenGen.AssertNotCalled(t, "GenerateToken", mock.Anything)
}

func TestLoginUserUseCase_Execute_InvalidPassword(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)

	useCase := NewLoginUserUseCase(mockRepo, mockTokenGen)

	ctx := context.Background()
	email := "test@example.com"
	password := "wrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	hashedPasswordStr := string(hashedPassword)

	existingUser := &user.User{
		ID:       uuid.New(),
		Email:    email,
		Password: &hashedPasswordStr,
	}

	mockRepo.On("FindByEmail", ctx, email).Return(existingUser, nil)

	// Act
	token, err := useCase.Execute(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid password")
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
	mockTokenGen.AssertNotCalled(t, "GenerateToken", mock.Anything)
}

func TestLoginUserUseCase_Execute_TokenGenerationError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)

	useCase := NewLoginUserUseCase(mockRepo, mockTokenGen)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedPasswordStr := string(hashedPassword)
	tokenError := errors.New("token generation failed")

	existingUser := &user.User{
		ID:       uuid.New(),
		Email:    email,
		Password: &hashedPasswordStr,
	}

	mockRepo.On("FindByEmail", ctx, email).Return(existingUser, nil)
	mockTokenGen.On("GenerateToken", existingUser).Return("", tokenError)

	// Act
	token, err := useCase.Execute(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
	mockTokenGen.AssertExpectations(t)
}
