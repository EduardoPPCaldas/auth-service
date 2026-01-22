package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EduardoPPCaldas/auth-service/internal/application/user/dto"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUseCases for testing handlers
type MockCreateUserUseCase struct {
	mock.Mock
}

func (m *MockCreateUserUseCase) Execute(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

type MockLoginUserUseCase struct {
	mock.Mock
}

func (m *MockLoginUserUseCase) Execute(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

type MockLoginWithGoogleUseCase struct {
	mock.Mock
}

func (m *MockLoginWithGoogleUseCase) Execute(ctx context.Context, idToken string) (string, error) {
	args := m.Called(ctx, idToken)
	return args.String(0), args.Error(1)
}

type MockGoogleOAuthService struct {
	mock.Mock
}

func (m *MockGoogleOAuthService) GetAuthURL() string {
	args := m.Called()
	return args.String(0)
}

func setupEcho() *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	return e
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func TestAuthHandler_CreateUser_Success(t *testing.T) {
	// Arrange
	mockCreateUser := new(MockCreateUserUseCase)
	mockLoginUser := new(MockLoginUserUseCase)
	mockLoginWithGoogle := new(MockLoginWithGoogleUseCase)
	mockGoogleOAuth := new(MockGoogleOAuthService)

	handler := NewAuthHandler(
		mockCreateUser,
		mockLoginUser,
		mockLoginWithGoogle,
		nil, // RefreshTokenUseCase - not needed for this test
		nil, // LogoutUseCase - not needed for this test
		mockGoogleOAuth,
	)

	req := dto.CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	token := "jwt-token-here"

	mockCreateUser.On("Execute", req.Email, req.Password).Return(token, nil)

	body, _ := json.Marshal(req)
	e := setupEcho()
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	// Act
	err := handler.CreateUser(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response dto.AuthResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, token, response.AccessToken)

	mockCreateUser.AssertExpectations(t)
}

func TestAuthHandler_CreateUser_InvalidRequest(t *testing.T) {
	// Arrange
	mockCreateUser := new(MockCreateUserUseCase)
	mockLoginUser := new(MockLoginUserUseCase)
	mockLoginWithGoogle := new(MockLoginWithGoogleUseCase)
	mockGoogleOAuth := new(MockGoogleOAuthService)

	handler := NewAuthHandler(
		mockCreateUser,
		mockLoginUser,
		mockLoginWithGoogle,
		nil, // RefreshTokenUseCase - not needed for this test
		nil, // LogoutUseCase - not needed for this test
		mockGoogleOAuth,
	)

	req := dto.CreateUserRequest{
		Email:    "invalid-email",
		Password: "short",
	}

	body, _ := json.Marshal(req)
	e := setupEcho()
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	// Act
	err := handler.CreateUser(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockCreateUser.AssertNotCalled(t, "Execute", mock.Anything, mock.Anything)
}

func TestAuthHandler_CreateUser_UseCaseError(t *testing.T) {
	// Arrange
	mockCreateUser := new(MockCreateUserUseCase)
	mockLoginUser := new(MockLoginUserUseCase)
	mockLoginWithGoogle := new(MockLoginWithGoogleUseCase)
	mockGoogleOAuth := new(MockGoogleOAuthService)

	handler := NewAuthHandler(
		mockCreateUser,
		mockLoginUser,
		mockLoginWithGoogle,
		nil, // RefreshTokenUseCase - not needed for this test
		nil, // LogoutUseCase - not needed for this test
		mockGoogleOAuth,
	)

	req := dto.CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockCreateUser.On("Execute", req.Email, req.Password).Return("", assert.AnError)

	body, _ := json.Marshal(req)
	e := setupEcho()
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	// Act
	err := handler.CreateUser(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockCreateUser.AssertExpectations(t)
}

func TestAuthHandler_LoginUser_Success(t *testing.T) {
	// Arrange
	mockCreateUser := new(MockCreateUserUseCase)
	mockLoginUser := new(MockLoginUserUseCase)
	mockLoginWithGoogle := new(MockLoginWithGoogleUseCase)
	mockGoogleOAuth := new(MockGoogleOAuthService)

	handler := NewAuthHandler(
		mockCreateUser,
		mockLoginUser,
		mockLoginWithGoogle,
		nil, // RefreshTokenUseCase - not needed for this test
		nil, // LogoutUseCase - not needed for this test
		mockGoogleOAuth,
	)

	req := dto.LoginUserRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	token := "jwt-token-here"

	mockLoginUser.On("Execute", req.Email, req.Password).Return(token, nil)

	body, _ := json.Marshal(req)
	e := setupEcho()
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	// Act
	err := handler.LoginUser(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.AuthResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, token, response.AccessToken)

	mockLoginUser.AssertExpectations(t)
}

func TestAuthHandler_LoginUser_InvalidCredentials(t *testing.T) {
	// Arrange
	mockCreateUser := new(MockCreateUserUseCase)
	mockLoginUser := new(MockLoginUserUseCase)
	mockLoginWithGoogle := new(MockLoginWithGoogleUseCase)
	mockGoogleOAuth := new(MockGoogleOAuthService)

	handler := NewAuthHandler(
		mockCreateUser,
		mockLoginUser,
		mockLoginWithGoogle,
		nil, // RefreshTokenUseCase - not needed for this test
		nil, // LogoutUseCase - not needed for this test
		mockGoogleOAuth,
	)

	req := dto.LoginUserRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	mockLoginUser.On("Execute", req.Email, req.Password).Return("", assert.AnError)

	body, _ := json.Marshal(req)
	e := setupEcho()
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	// Act
	err := handler.LoginUser(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	mockLoginUser.AssertExpectations(t)
}

func TestAuthHandler_ChallengeGoogleAuth_Success(t *testing.T) {
	// Arrange
	mockCreateUser := new(MockCreateUserUseCase)
	mockLoginUser := new(MockLoginUserUseCase)
	mockLoginWithGoogle := new(MockLoginWithGoogleUseCase)
	mockGoogleOAuth := new(MockGoogleOAuthService)

	handler := NewAuthHandler(
		mockCreateUser,
		mockLoginUser,
		mockLoginWithGoogle,
		nil, // RefreshTokenUseCase - not needed for this test
		nil, // LogoutUseCase - not needed for this test
		mockGoogleOAuth,
	)

	authURL := "https://accounts.google.com/oauth/authorize"
	mockGoogleOAuth.On("GetAuthURL").Return(authURL)

	e := setupEcho()
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/challenge", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	// Act
	err := handler.ChallengeGoogleAuth(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, authURL, rec.Header().Get("Location"))

	mockGoogleOAuth.AssertExpectations(t)
}
