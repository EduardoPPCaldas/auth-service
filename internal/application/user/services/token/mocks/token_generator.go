package mocks

import (
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/google/uuid"
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

func (m *MockTokenGenerator) ExtractUserID(tokenString string) (uuid.UUID, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return uuid.Nil, args.Error(1)
	}
	return args.Get(0).(uuid.UUID), args.Error(1)
}
