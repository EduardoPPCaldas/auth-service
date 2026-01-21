package mocks

import (
	"context"

	"github.com/EduardoPPCaldas/auth-service/internal/application/user/services/oauth"
	"github.com/stretchr/testify/mock"
)

// MockGoogleTokenValidator is a mock implementation of GoogleTokenValidator
type MockGoogleTokenValidator struct {
	mock.Mock
}

func (m *MockGoogleTokenValidator) Validate(ctx context.Context, idToken string) (*oauth.GoogleUser, error) {
	args := m.Called(ctx, idToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oauth.GoogleUser), args.Error(1)
}
