package mocks

import (
	"context"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/role"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) FindByID(id uuid.UUID) (*role.Role, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*role.Role), args.Error(1)
}

func (m *MockRoleRepository) FindByName(name string) (*role.Role, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*role.Role), args.Error(1)
}

func (m *MockRoleRepository) Create(ctx context.Context, r *role.Role) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockRoleRepository) FindOrCreateDefault() (*role.Role, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*role.Role), args.Error(1)
}

func (m *MockRoleRepository) IsRBACEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}
