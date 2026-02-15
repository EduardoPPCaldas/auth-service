package usecases

import (
	"context"
	"errors"
	"testing"

	tokenmocks "github.com/EduardoPPCaldas/auth-service/internal/application/user/services/token/mocks"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/role"
	rolemocks "github.com/EduardoPPCaldas/auth-service/internal/domain/role/mocks"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateUserUseCase_Execute_Success_WithRBAC(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockRoleRepo := new(rolemocks.MockRoleRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)

	useCase := NewCreateUserUseCase(mockRepo, mockRoleRepo, mockTokenGen)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	expectedToken := "jwt-token-here"
	defaultRole := role.NewUserRole()

	mockRepo.On("FindByEmail", email).Return(nil, gorm.ErrRecordNotFound)
	mockRoleRepo.On("IsRBACEnabled", ctx).Return(true)
	mockRoleRepo.On("FindOrCreateDefault", ctx).Return(defaultRole, nil)
	mockRepo.On("Create", mock.AnythingOfType("*user.User")).Return(nil)
	mockTokenGen.On("GenerateToken", mock.AnythingOfType("*user.User")).Return(expectedToken, nil)

	// Act
	token, err := useCase.Execute(ctx, email, password)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
	mockRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
	mockTokenGen.AssertExpectations(t)
}

func TestCreateUserUseCase_Execute_Success_WithoutRBAC(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockRoleRepo := new(rolemocks.MockRoleRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)

	useCase := NewCreateUserUseCase(mockRepo, mockRoleRepo, mockTokenGen)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	expectedToken := "jwt-token-here"

	mockRepo.On("FindByEmail", email).Return(nil, gorm.ErrRecordNotFound)
	mockRoleRepo.On("IsRBACEnabled", ctx).Return(false)
	mockRepo.On("Create", mock.AnythingOfType("*user.User")).Return(nil)
	mockTokenGen.On("GenerateToken", mock.AnythingOfType("*user.User")).Return(expectedToken, nil)

	// Act
	token, err := useCase.Execute(ctx, email, password)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
	mockRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
	mockTokenGen.AssertExpectations(t)
}

func TestCreateUserUseCase_Execute_UserAlreadyExists(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockRoleRepo := new(rolemocks.MockRoleRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)

	useCase := NewCreateUserUseCase(mockRepo, mockRoleRepo, mockTokenGen)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	existingUser := &user.User{
		ID:    uuid.New(),
		Email: email,
	}

	mockRepo.On("FindByEmail", email).Return(existingUser, nil)

	// Act
	token, err := useCase.Execute(ctx, email, password)

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
	mockRoleRepo := new(rolemocks.MockRoleRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)

	useCase := NewCreateUserUseCase(mockRepo, mockRoleRepo, mockTokenGen)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	repoError := errors.New("database connection failed")

	mockRepo.On("FindByEmail", email).Return(nil, repoError)

	// Act
	token, err := useCase.Execute(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating user")
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
	mockTokenGen.AssertNotCalled(t, "GenerateToken", mock.Anything)
}

func TestCreateUserUseCase_Execute_CreateError_WithRBAC(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockRoleRepo := new(rolemocks.MockRoleRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)

	useCase := NewCreateUserUseCase(mockRepo, mockRoleRepo, mockTokenGen)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	createError := errors.New("failed to create user")
	defaultRole := role.NewUserRole()

	mockRepo.On("FindByEmail", email).Return(nil, gorm.ErrRecordNotFound)
	mockRoleRepo.On("IsRBACEnabled", ctx).Return(true)
	mockRoleRepo.On("FindOrCreateDefault", ctx).Return(defaultRole, nil)
	mockRepo.On("Create", mock.AnythingOfType("*user.User")).Return(createError)

	// Act
	token, err := useCase.Execute(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating user")
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
	mockTokenGen.AssertNotCalled(t, "GenerateToken", mock.Anything)
}

func TestCreateUserUseCase_Execute_TokenGenerationError_WithRBAC(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockRoleRepo := new(rolemocks.MockRoleRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)

	useCase := NewCreateUserUseCase(mockRepo, mockRoleRepo, mockTokenGen)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	tokenError := errors.New("token generation failed")
	defaultRole := role.NewUserRole()

	mockRepo.On("FindByEmail", email).Return(nil, gorm.ErrRecordNotFound)
	mockRoleRepo.On("IsRBACEnabled", ctx).Return(true)
	mockRoleRepo.On("FindOrCreateDefault", ctx).Return(defaultRole, nil)
	mockRepo.On("Create", mock.AnythingOfType("*user.User")).Return(nil)
	mockTokenGen.On("GenerateToken", mock.AnythingOfType("*user.User")).Return("", tokenError)

	// Act
	token, err := useCase.Execute(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
	mockTokenGen.AssertExpectations(t)
}
