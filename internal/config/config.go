package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	// Database
	DatabaseURL string

	// Server
	Port string

	// JWT
	JWTSecret     string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration

	// Google OAuth
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURI  string
	GoogleCallbackURL  string
	OAuthState         string
}

func Load() *Config {
	// Database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/authdb?sslmode=disable"
	}

	// Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable is required")
	}

	accessExpiry, _ := time.ParseDuration(getEnvOrDefault("JWT_ACCESS_EXPIRY", "24h"))
	refreshExpiry, _ := time.ParseDuration(getEnvOrDefault("JWT_REFRESH_EXPIRY", "168h"))

	// JWT Refresh Secret
	jwtRefreshSecret := os.Getenv("JWT_REFRESH_SECRET")

	// Google OAuth
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURI := os.Getenv("GOOGLE_REDIRECT_URI")

	// Additional OAuth configuration
	googleCallbackURL := os.Getenv("GOOGLE_CALLBACK_URL")
	oauthState := os.Getenv("OAUTH_STATE")

	return &Config{
		DatabaseURL:         dbURL,
		Port:              port,
		JWTSecret:         jwtSecret,
		JWTRefreshSecret:    jwtRefreshSecret,
		AccessExpiry:        accessExpiry,
		RefreshExpiry:        refreshExpiry,
		GoogleClientID:        googleClientID,
		GoogleClientSecret:     googleClientSecret,
		GoogleRedirectURI:      googleRedirectURI,
		GoogleCallbackURL:     googleCallbackURL,
		OAuthState:          oauthState,
	}
}

	// Server
	if port == "" {
		port = "8080"
	}

	// JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable is required")
	}

	accessExpiry, _ := time.ParseDuration(getEnvOrDefault("JWT_ACCESS_EXPIRY", "24h"))
	refreshExpiry, _ := time.ParseDuration(getEnvOrDefault("JWT_REFRESH_EXPIRY", "168h"))

	// JWT Refresh Secret
	jwtRefreshSecret := os.Getenv("JWT_REFRESH_SECRET")

	// Google OAuth
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURI := os.Getenv("GOOGLE_REDIRECT_URI")

	// Additional OAuth configuration (from existing patterns)
	googleCallbackURL := os.Getenv("GOOGLE_CALLBACK_URL")
	oauthState := os.Getenv("OAUTH_STATE")

	return &Config{
		DatabaseURL:        dbURL,
		Port:               port,
		JWTSecret:          jwtSecret,
		JWTRefreshSecret:   jwtRefreshSecret,
		AccessExpiry:       accessExpiry,
		RefreshExpiry:      refreshExpiry,
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
		GoogleRedirectURI:  googleRedirectURI,
		GoogleCallbackURL:  googleCallbackURL,
		OAuthState:         oauthState,
	}

	// Server

	port := os.Getenv("PORT")

	// JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable is required")
	}

	accessExpiry, _ := time.ParseDuration(getEnvOrDefault("JWT_ACCESS_EXPIRY", "24h"))
	refreshExpiry, _ := time.ParseDuration(getEnvOrDefault("JWT_REFRESH_EXPIRY", "168h"))

	// Google OAuth
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURI := os.Getenv("GOOGLE_REDIRECT_URI")

	return &Config{
		DatabaseURL:        dbURL,
		Port:               port,
		JWTSecret:          jwtSecret,
		AccessExpiry:       accessExpiry,
		RefreshExpiry:      refreshExpiry,
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
		GoogleRedirectURI:  googleRedirectURI,
	}
}

func (c *Config) Validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

func (c *Config) SetupOAuthDefaults() {
	if c.GoogleRedirectURI == "" {
		c.GoogleRedirectURI = fmt.Sprintf("http://localhost:%s/api/v1/auth/google/callback", c.Port)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
