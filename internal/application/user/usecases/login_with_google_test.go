package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/EduardoPPCaldas/auth-service/internal/application/user/services/oauth"
	oauthmocks "github.com/EduardoPPCaldas/auth-service/internal/application/user/services/oauth/mocks"
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

func TestLoginWithGoogleUseCase_Execute_NewUser_WithRBAC(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockRoleRepo := new(rolemocks.MockRoleRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)
	mockGoogleValidator := new(oauthmocks.MockGoogleTokenValidator)

	useCase := NewLoginWithGoogleUseCase(mockRepo, mockRoleRepo, mockTokenGen, mockGoogleValidator)

	ctx := context.Background()
	idToken := "google-id-token"
	email := "google@example.com"
	expectedToken := "jwt-token-here"
	defaultRole := role.NewUserRole()

	googleUser := &oauth.GoogleUser{
		Email: email,
		Name:  "Google User",
	}

	mockGoogleValidator.On("Validate", ctx, idToken).Return(googleUser, nil)
	mockRepo.On("FindByEmail", email).Return(nil, gorm.ErrRecordNotFound)
	mockRoleRepo.On("IsRBACEnabled").Return(true)
	mockRoleRepo.On("FindOrCreateDefault").Return(defaultRole, nil)
	mockRepo.On("Create", mock.AnythingOfType("*user.User")).Return(nil)
	mockTokenGen.On("GenerateToken", mock.AnythingOfType("*user.User")).Return(expectedToken, nil)

	// Act
	token, err := useCase.Execute(ctx, idToken)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
	mockGoogleValidator.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
	mockTokenGen.AssertExpectations(t)
}

func TestLoginWithGoogleUseCase_Execute_NewUser_WithoutRBAC(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockRoleRepo := new(rolemocks.MockRoleRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)
	mockGoogleValidator := new(oauthmocks.MockGoogleTokenValidator)

	useCase := NewLoginWithGoogleUseCase(mockRepo, mockRoleRepo, mockTokenGen, mockGoogleValidator)

	ctx := context.Background()
	idToken := "google-id-token"
	email := "google@example.com"
	expectedToken := "jwt-token-here"

	googleUser := &oauth.GoogleUser{
		Email: email,
		Name:  "Google User",
	}

	mockGoogleValidator.On("Validate", ctx, idToken).Return(googleUser, nil)
	mockRepo.On("FindByEmail", email).Return(nil, gorm.ErrRecordNotFound)
	mockRoleRepo.On("IsRBACEnabled").Return(false)
	mockRepo.On("Create", mock.AnythingOfType("*user.User")).Return(nil)
	mockTokenGen.On("GenerateToken", mock.AnythingOfType("*user.User")).Return(expectedToken, nil)

	// Act
	token, err := useCase.Execute(ctx, idToken)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
	mockGoogleValidator.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
	mockTokenGen.AssertExpectations(t)
}

func TestLoginWithGoogleUseCase_Execute_ExistingUser(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockRoleRepo := new(rolemocks.MockRoleRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)
	mockGoogleValidator := new(oauthmocks.MockGoogleTokenValidator)

	useCase := NewLoginWithGoogleUseCase(mockRepo, mockRoleRepo, mockTokenGen, mockGoogleValidator)

	ctx := context.Background()
	idToken := "google-id-token"
	email := "google@example.com"
	expectedToken := "jwt-token-here"

	googleUser := &oauth.GoogleUser{
		Email: email,
		Name:  "Google User",
	}

	existingUser := &user.User{
		ID:    uuid.New(),
		Email: email,
	}

	mockGoogleValidator.On("Validate", ctx, idToken).Return(googleUser, nil)
	mockRepo.On("FindByEmail", email).Return(existingUser, nil)
	mockTokenGen.On("GenerateToken", existingUser).Return(expectedToken, nil)

	// Act
	token, err := useCase.Execute(ctx, idToken)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
	mockGoogleValidator.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockTokenGen.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestLoginWithGoogleUseCase_Execute_InvalidToken(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockRoleRepo := new(rolemocks.MockRoleRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)
	mockGoogleValidator := new(oauthmocks.MockGoogleTokenValidator)

	useCase := NewLoginWithGoogleUseCase(mockRepo, mockRoleRepo, mockTokenGen, mockGoogleValidator)

	ctx := context.Background()
	idToken := "invalid-token"
	validationError := errors.New("invalid token")

	mockGoogleValidator.On("Validate", ctx, idToken).Return(nil, validationError)

	// Act
	token, err := useCase.Execute(ctx, idToken)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid google token")
	assert.Empty(t, token)
	mockGoogleValidator.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "FindByEmail", mock.Anything)
	mockTokenGen.AssertNotCalled(t, "GenerateToken", mock.Anything)
}

func TestLoginWithGoogleUseCase_Execute_CreateError_WithRBAC(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockRoleRepo := new(rolemocks.MockRoleRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)
	mockGoogleValidator := new(oauthmocks.MockGoogleTokenValidator)

	useCase := NewLoginWithGoogleUseCase(mockRepo, mockRoleRepo, mockTokenGen, mockGoogleValidator)

	ctx := context.Background()
	idToken := "google-id-token"
	email := "google@example.com"
	createError := errors.New("failed to create user")
	defaultRole := role.NewUserRole()

	googleUser := &oauth.GoogleUser{
		Email: email,
		Name:  "Google User",
	}

	mockGoogleValidator.On("Validate", ctx, idToken).Return(googleUser, nil)
	mockRepo.On("FindByEmail", email).Return(nil, gorm.ErrRecordNotFound)
	mockRoleRepo.On("IsRBACEnabled").Return(true)
	mockRoleRepo.On("FindOrCreateDefault").Return(defaultRole, nil)
	mockRepo.On("Create", mock.AnythingOfType("*user.User")).Return(createError)

	// Act
	token, err := useCase.Execute(ctx, idToken)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create user")
	assert.Empty(t, token)
	mockGoogleValidator.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
	mockTokenGen.AssertNotCalled(t, "GenerateToken", mock.Anything)
}

func TestLoginWithGoogleUseCase_Execute_FindByEmailError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockRoleRepo := new(rolemocks.MockRoleRepository)
	mockTokenGen := new(tokenmocks.MockTokenGenerator)
	mockGoogleValidator := new(oauthmocks.MockGoogleTokenValidator)

	useCase := NewLoginWithGoogleUseCase(mockRepo, mockRoleRepo, mockTokenGen, mockGoogleValidator)

	ctx := context.Background()
	idToken := "google-id-token"
	email := "google@example.com"
	findError := errors.New("database error")

	googleUser := &oauth.GoogleUser{
		Email: email,
		Name:  "Google User",
	}

	mockGoogleValidator.On("Validate", ctx, idToken).Return(googleUser, nil)
	mockRepo.On("FindByEmail", email).Return(nil, findError)

	// Act
	token, err := useCase.Execute(ctx, idToken)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find user")
	assert.Empty(t, token)
	mockGoogleValidator.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockTokenGen.AssertNotCalled(t, "GenerateToken", mock.Anything)
}
