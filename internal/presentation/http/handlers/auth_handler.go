package handlers

import (
	"context"
	"net/http"

	"github.com/EduardoPPCaldas/auth-service/internal/application/user/dto"
	"github.com/EduardoPPCaldas/auth-service/internal/application/user/usecases"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	createUserUseCase      CreateUserUseCase
	loginUserUseCase       LoginUserUseCase
	loginWithGoogleUseCase LoginWithGoogleUseCase
	refreshTokenUseCase    RefreshTokenUseCase
	logoutUseCase          LogoutUseCase
	googleOAuthService     GoogleOAuthService
}

type CreateUserUseCase interface {
	Execute(ctx context.Context, email, password string) (string, error)
}

type LoginUserUseCase interface {
	Execute(email, password string) (string, error)
}

type LoginWithGoogleUseCase interface {
	Execute(ctx context.Context, idToken string) (string, error)
}

type RefreshTokenUseCase interface {
	Execute(ctx context.Context, refreshToken string) (*usecases.RefreshTokenResponse, error)
}

type LogoutUseCase interface {
	Execute(ctx context.Context, userID string) error
	LogoutSingle(ctx context.Context, refreshToken string) error
}

type GoogleOAuthService interface {
	GetAuthURL() string
}

func NewAuthHandler(
	createUserUseCase CreateUserUseCase,
	loginUserUseCase LoginUserUseCase,
	loginWithGoogleUseCase LoginWithGoogleUseCase,
	refreshTokenUseCase RefreshTokenUseCase,
	logoutUseCase LogoutUseCase,
	googleOAuthService GoogleOAuthService,
) *AuthHandler {
	return &AuthHandler{
		createUserUseCase:      createUserUseCase,
		loginUserUseCase:       loginUserUseCase,
		loginWithGoogleUseCase: loginWithGoogleUseCase,
		refreshTokenUseCase:    refreshTokenUseCase,
		logoutUseCase:          logoutUseCase,
		googleOAuthService:     googleOAuthService,
	}
}

// CreateUser handles user registration
// POST /api/v1/auth/register
func (h *AuthHandler) CreateUser(c echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	token, err := h.createUserUseCase.Execute(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.AuthResponse{AccessToken: token, RefreshToken: "", TokenType: "Bearer", ExpiresIn: 86400})
}

// LoginUser handles user login
// POST /api/v1/auth/login
func (h *AuthHandler) LoginUser(c echo.Context) error {
	var req dto.LoginUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	token, err := h.loginUserUseCase.Execute(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.AuthResponse{AccessToken: token, RefreshToken: "", TokenType: "Bearer", ExpiresIn: 86400})
}

// LoginWithGoogle handles Google OAuth login
// POST /api/v1/auth/login/google
func (h *AuthHandler) LoginWithGoogle(c echo.Context) error {
	var req dto.LoginWithGoogleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	token, err := h.loginWithGoogleUseCase.Execute(c.Request().Context(), req.IDToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.AuthResponse{AccessToken: token, RefreshToken: "", TokenType: "Bearer", ExpiresIn: 86400})
}

// RefreshToken handles token refresh
// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req dto.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	response, err := h.refreshTokenUseCase.Execute(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}

// Logout handles user logout (revokes a specific refresh token)
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c echo.Context) error {
	var req dto.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := h.logoutUseCase.LogoutSingle(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "logged out successfully"})
}

// LogoutAll handles logout from all devices (revokes all user refresh tokens)
// POST /api/v1/auth/logout-all
func (h *AuthHandler) LogoutAll(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not authenticated"})
	}

	userIDStr, ok := userID.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user ID format"})
	}

	err := h.logoutUseCase.Execute(c.Request().Context(), userIDStr.String())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "logged out from all devices successfully"})
}

// ChallengeGoogleAuth handles Google OAuth challenge (redirects to OAuth URL)
// GET /api/v1/auth/google/challenge
func (h *AuthHandler) ChallengeGoogleAuth(c echo.Context) error {
	authURL := h.googleOAuthService.GetAuthURL()
	return c.Redirect(http.StatusFound, authURL)
}
