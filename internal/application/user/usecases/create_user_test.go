package usecases

import (
	"errors"
	"testing"

	tokenmocks "github.com/EduardoPPCaldas/auth-service/internal/application/user/services/token/mocks"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateUserUseCase_Execute_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)

	useCase := NewCreateUserUseCase(mockRepo, mockTokenGen)

	email := "test@example.com"
	password := "password123"
	expectedToken := "jwt-token-here"

	mockRepo.On("FindByEmail", email).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Create", mock.AnythingOfType("*user.User")).Return(nil)
	mockTokenGen.On("GenerateToken", mock.AnythingOfType("*user.User")).Return(expectedToken, nil)

	// Act
	token, err := useCase.Execute(email, password)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
	mockRepo.AssertExpectations(t)
	mockTokenGen.AssertExpectations(t)
}

func TestCreateUserUseCase_Execute_UserAlreadyExists(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)

	useCase := NewCreateUserUseCase(mockRepo, mockTokenGen)

	email := "test@example.com"
	password := "password123"
	existingUser := &user.User{
		ID:    uuid.New(),
		Email: email,
	}

	mockRepo.On("FindByEmail", email).Return(existingUser, nil)

	// Act
	token, err := useCase.Execute(email, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "user already exists", err.Error())
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
	mockTokenGen.AssertNotCalled(t, "GenerateToken", mock.Anything)
}

func TestCreateUserUseCase_Execute_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)

	useCase := NewCreateUserUseCase(mockRepo, mockTokenGen)

	email := "test@example.com"
	password := "password123"
	repoError := errors.New("database connection failed")

	mockRepo.On("FindByEmail", email).Return(nil, repoError)

	// Act
	token, err := useCase.Execute(email, password)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating user")
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
	mockTokenGen.AssertNotCalled(t, "GenerateToken", mock.Anything)
}

func TestCreateUserUseCase_Execute_CreateError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)

	useCase := NewCreateUserUseCase(mockRepo, mockTokenGen)

	email := "test@example.com"
	password := "password123"
	createError := errors.New("failed to create user")

	mockRepo.On("FindByEmail", email).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Create", mock.AnythingOfType("*user.User")).Return(createError)

	// Act
	token, err := useCase.Execute(email, password)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating user")
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
	mockTokenGen.AssertNotCalled(t, "GenerateToken", mock.Anything)
}

func TestCreateUserUseCase_Execute_TokenGenerationError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)

	useCase := NewCreateUserUseCase(mockRepo, mockTokenGen)

	email := "test@example.com"
	password := "password123"
	tokenError := errors.New("token generation failed")

	mockRepo.On("FindByEmail", email).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Create", mock.AnythingOfType("*user.User")).Return(nil)
	mockTokenGen.On("GenerateToken", mock.AnythingOfType("*user.User")).Return("", tokenError)

	// Act
	token, err := useCase.Execute(email, password)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
	mockTokenGen.AssertExpectations(t)
}
