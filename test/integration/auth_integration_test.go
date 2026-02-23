package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/EduardoPPCaldas/auth-service/internal/application/user/dto"
	"github.com/EduardoPPCaldas/auth-service/internal/application/user/services/token"
	"github.com/EduardoPPCaldas/auth-service/internal/application/user/usecases"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/role"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/EduardoPPCaldas/auth-service/internal/infrastructure/oauth/google"
	postgresRepo "github.com/EduardoPPCaldas/auth-service/internal/infrastructure/postgres/repository"
	httphandler "github.com/EduardoPPCaldas/auth-service/internal/presentation/http"
	"github.com/EduardoPPCaldas/auth-service/internal/presentation/http/handlers"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupIntegrationServer(t *testing.T) (*echo.Echo, func()) {
	ctx := context.Background()

	postgresContainer, err := postgrescontainer.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgrescontainer.WithDatabase("testdb"),
		postgrescontainer.WithUsername("testuser"),
		postgrescontainer.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err)

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&user.User{}, &role.Role{}, &role.Permission{})
	require.NoError(t, err)

	// Set JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key-for-integration-tests")

	// Initialize repositories
	userRepo := postgresRepo.NewUserRepository(db)
	roleRepo := postgresRepo.NewRoleRepository(db)

	// Initialize services
	tokenGenerator := token.NewTokenGenerator()
	googleValidator := google.NewGoogleTokenValidator("")
	googleOAuthService := google.NewGoogleOAuthChallengeService("", "", "")

	// Initialize use cases
	createUserUseCase := usecases.NewCreateUserUseCase(userRepo, roleRepo, tokenGenerator)
	loginUserUseCase := usecases.NewLoginUserUseCase(userRepo, tokenGenerator)
	loginWithGoogleUseCase := usecases.NewLoginWithGoogleUseCase(userRepo, roleRepo, tokenGenerator, googleValidator)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(
		createUserUseCase,
		loginUserUseCase,
		loginWithGoogleUseCase,
		nil, // RefreshTokenUseCase - not needed for these integration tests
		nil, // LogoutUseCase - not needed for these integration tests
		googleOAuthService,
	)

	// Initialize Echo
	e := echo.New()
	httphandler.SetupMiddleware(e)
	e.Validator = &CustomValidator{validator: validator.New()}
	httphandler.SetupRoutes(e, authHandler, nil, nil, nil)

	cleanup := func() {
		os.Unsetenv("JWT_SECRET")
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}

	return e, cleanup
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func TestAuthIntegration_Register_Success(t *testing.T) {
	// Arrange
	e, cleanup := setupIntegrationServer(t)
	defer cleanup()

	req := dto.CreateUserRequest{
		Email:    "integration@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, httpReq)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response dto.AuthResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotEmpty(t, response.AccessToken)
}

func TestAuthIntegration_Register_Then_Login(t *testing.T) {
	// Arrange
	e, cleanup := setupIntegrationServer(t)
	defer cleanup()

	email := "registerlogin@example.com"
	password := "password123"

	// Register
	registerReq := dto.CreateUserRequest{
		Email:    email,
		Password: password,
	}

	registerBody, _ := json.Marshal(registerReq)
	registerHTTPReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(registerBody))
	registerHTTPReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	registerRec := httptest.NewRecorder()

	e.ServeHTTP(registerRec, registerHTTPReq)
	require.Equal(t, http.StatusCreated, registerRec.Code)

	var registerResponse dto.AuthResponse
	err := json.Unmarshal(registerRec.Body.Bytes(), &registerResponse)
	require.NoError(t, err)
	require.NotEmpty(t, registerResponse.AccessToken)

	// Login
	loginReq := dto.LoginUserRequest{
		Email:    email,
		Password: password,
	}

	loginBody, _ := json.Marshal(loginReq)
	loginHTTPReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	loginHTTPReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	loginRec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(loginRec, loginHTTPReq)

	// Assert
	assert.Equal(t, http.StatusOK, loginRec.Code)

	var loginResponse dto.AuthResponse
	err = json.Unmarshal(loginRec.Body.Bytes(), &loginResponse)
	require.NoError(t, err)
	assert.NotEmpty(t, loginResponse.AccessToken)
	// Tokens may be the same if generated at the same time with the same expiration
	// What's important is that login succeeds and returns a valid token
}

func TestAuthIntegration_Register_DuplicateEmail(t *testing.T) {
	// Arrange
	e, cleanup := setupIntegrationServer(t)
	defer cleanup()

	email := "duplicate@example.com"
	password := "password123"

	// First registration
	req1 := dto.CreateUserRequest{
		Email:    email,
		Password: password,
	}

	body1, _ := json.Marshal(req1)
	httpReq1 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body1))
	httpReq1.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec1 := httptest.NewRecorder()

	e.ServeHTTP(rec1, httpReq1)
	require.Equal(t, http.StatusCreated, rec1.Code)

	// Second registration with same email
	req2 := dto.CreateUserRequest{
		Email:    email,
		Password: "differentpassword",
	}

	body2, _ := json.Marshal(req2)
	httpReq2 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body2))
	httpReq2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec2 := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec2, httpReq2)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rec2.Code)
}

func TestAuthIntegration_Login_InvalidCredentials(t *testing.T) {
	// Arrange
	e, cleanup := setupIntegrationServer(t)
	defer cleanup()

	// Register first
	registerReq := dto.CreateUserRequest{
		Email:    "invalidlogin@example.com",
		Password: "correctpassword",
	}

	registerBody, _ := json.Marshal(registerReq)
	registerHTTPReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(registerBody))
	registerHTTPReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	registerRec := httptest.NewRecorder()

	e.ServeHTTP(registerRec, registerHTTPReq)
	require.Equal(t, http.StatusCreated, registerRec.Code)

	// Try to login with wrong password
	loginReq := dto.LoginUserRequest{
		Email:    "invalidlogin@example.com",
		Password: "wrongpassword",
	}

	loginBody, _ := json.Marshal(loginReq)
	loginHTTPReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	loginHTTPReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	loginRec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(loginRec, loginHTTPReq)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, loginRec.Code)
}

func TestAuthIntegration_Register_InvalidRequest(t *testing.T) {
	// Arrange
	e, cleanup := setupIntegrationServer(t)
	defer cleanup()

	req := dto.CreateUserRequest{
		Email:    "invalid-email",
		Password: "short",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, httpReq)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
