package handlers

import (
	"net/http"

	"github.com/EduardoPPCaldas/auth-service/internal/application/user/dto"
	"github.com/EduardoPPCaldas/auth-service/internal/application/user/usecases"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	createUserUseCase      *usecases.CreateUserUseCase
	loginUserUseCase       *usecases.LoginUserUseCase
	loginWithGoogleUseCase *usecases.LoginWithGoogleUseCase
	googleOAuthService     GoogleOAuthService
}

type GoogleOAuthService interface {
	GetAuthURL() string
}

func NewAuthHandler(
	createUserUseCase *usecases.CreateUserUseCase,
	loginUserUseCase *usecases.LoginUserUseCase,
	loginWithGoogleUseCase *usecases.LoginWithGoogleUseCase,
	googleOAuthService GoogleOAuthService,
) *AuthHandler {
	return &AuthHandler{
		createUserUseCase:      createUserUseCase,
		loginUserUseCase:       loginUserUseCase,
		loginWithGoogleUseCase: loginWithGoogleUseCase,
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

	token, err := h.createUserUseCase.Execute(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.AuthResponse{Token: token})
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

	return c.JSON(http.StatusOK, dto.AuthResponse{Token: token})
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

	return c.JSON(http.StatusOK, dto.AuthResponse{Token: token})
}

// ChallengeGoogleAuth handles Google OAuth challenge (returns OAuth URL)
// GET /api/v1/auth/google/challenge
func (h *AuthHandler) ChallengeGoogleAuth(c echo.Context) error {
	authURL := h.googleOAuthService.GetAuthURL()
	return c.JSON(http.StatusOK, dto.GoogleOAuthChallengeResponse{AuthURL: authURL})
}
