package dto

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
	Token string `json:"token"`
}

// GoogleOAuthChallengeResponse represents the response for Google OAuth challenge
type GoogleOAuthChallengeResponse struct {
	AuthURL string `json:"auth_url"`
}
