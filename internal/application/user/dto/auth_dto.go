package dto

import (
	"time"
)

// CreateUserRequest represents the request body for user registration
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginUserRequest represents the request body for user login
type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginWithGoogleRequest represents the request body for Google OAuth login
type LoginWithGoogleRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}

// AuthResponse represents the response for authentication endpoints
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// RefreshTokenRequest represents the request body for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse represents the response for token refresh
type RefreshTokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// GoogleOAuthChallengeResponse represents the response for Google OAuth challenge
type GoogleOAuthChallengeResponse struct {
	AuthURL string `json:"auth_url"`
}
