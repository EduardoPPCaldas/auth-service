package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/role"
	rolemocks "github.com/EduardoPPCaldas/auth-service/internal/domain/role/mocks"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	usermocks "github.com/EduardoPPCaldas/auth-service/internal/domain/user/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateRoleUseCase_Execute_Success(t *testing.T) {
	mockRoleRepo := new(rolemocks.MockRoleRepository)
	mockUserRepo := new(usermocks.MockUserRepository)
	useCase := NewCreateRoleUseCase(mockRoleRepo, mockUserRepo)

	adminUser := &user.User{
		ID:   uuid.New(),
		Role: role.NewAdminRole(),
	}

	ctx := context.Background()
	input := CreateRoleInput{
		AdminUserID: adminUser.ID,
		Name:        "moderator",
		Permissions: []string{"posts:read", "posts:write"},
	}

	mockUserRepo.On("FindByID", ctx, adminUser.ID).Return(adminUser, nil)
	mockRoleRepo.On("FindByName", ctx, input.Name).Return(nil, gorm.ErrRecordNotFound)
	mockRoleRepo.On("Create", ctx, mock.AnythingOfType("*role.Role")).Return(nil)

	result, err := useCase.Execute(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, input.Name, result.Name)
	assert.Equal(t, len(input.Permissions), len(result.Permissions))
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestCreateRoleUseCase_Execute_NotAdmin(t *testing.T) {
	mockRoleRepo := new(rolemocks.MockRoleRepository)
	mockUserRepo := new(usermocks.MockUserRepository)
	useCase := NewCreateRoleUseCase(mockRoleRepo, mockUserRepo)

	regularUser := &user.User{
		ID:   uuid.New(),
		Role: role.NewUserRole(),
	}

	ctx := context.Background()
	input := CreateRoleInput{
		AdminUserID: regularUser.ID,
		Name:        "moderator",
		Permissions: []string{"posts:read"},
	}

	mockUserRepo.On("FindByID", ctx, regularUser.ID).Return(regularUser, nil)

	result, err := useCase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "admin privileges")
	mockUserRepo.AssertExpectations(t)
}

func TestCreateRoleUseCase_Execute_RoleAlreadyExists(t *testing.T) {
	mockRoleRepo := new(rolemocks.MockRoleRepository)
	mockUserRepo := new(usermocks.MockUserRepository)
	useCase := NewCreateRoleUseCase(mockRoleRepo, mockUserRepo)

	adminUser := &user.User{
		ID:   uuid.New(),
		Role: role.NewAdminRole(),
	}

	existingRole := role.New("moderator", []string{"posts:read"})

	ctx := context.Background()
	input := CreateRoleInput{
		AdminUserID: adminUser.ID,
		Name:        "moderator",
		Permissions: []string{"posts:read", "posts:write"},
	}

	mockUserRepo.On("FindByID", ctx, adminUser.ID).Return(adminUser, nil)
	mockRoleRepo.On("FindByName", ctx, input.Name).Return(existingRole, nil)

	result, err := useCase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists")
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestCreateRoleUseCase_Execute_AdminUserNotFound(t *testing.T) {
	mockRoleRepo := new(rolemocks.MockRoleRepository)
	mockUserRepo := new(usermocks.MockUserRepository)
	useCase := NewCreateRoleUseCase(mockRoleRepo, mockUserRepo)

	adminUserID := uuid.New()

	ctx := context.Background()
	input := CreateRoleInput{
		AdminUserID: adminUserID,
		Name:        "moderator",
		Permissions: []string{"posts:read"},
	}

	mockUserRepo.On("FindByID", ctx, adminUserID).Return(nil, errors.New("user not found"))

	result, err := useCase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "admin user not found")
	mockUserRepo.AssertExpectations(t)
}
