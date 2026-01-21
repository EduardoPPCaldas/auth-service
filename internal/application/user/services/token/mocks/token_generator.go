package mocks

import (
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/stretchr/testify/mock"
)

// MockTokenGenerator is a mock implementation of TokenGenerator
type MockTokenGenerator struct {
	mock.Mock
}

func (m *MockTokenGenerator) GenerateToken(u *user.User) (string, error) {
	args := m.Called(u)
	return args.String(0), args.Error(1)
}
