package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrTokenExpired     = errors.New("token has expired")
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrTokenMissing     = errors.New("authorization token is missing")
	ErrTokenMalformed   = errors.New("token is malformed")
	ErrPermissionDenied = errors.New("permission denied")
	ErrRoleDenied       = errors.New("role denied")
	ErrUserNotFound     = errors.New("user not found in context")
)

type AuthError struct {
	Type    AuthErrorType
	Message string
	Code    int
}

type AuthErrorType string

const (
	ErrorTypeExpired      AuthErrorType = "expired"
	ErrorTypeInvalid      AuthErrorType = "invalid"
	ErrorTypeMissing      AuthErrorType = "missing"
	ErrorTypeMalformed    AuthErrorType = "malformed"
	ErrorTypePermission   AuthErrorType = "permission_denied"
	ErrorTypeRole         AuthErrorType = "role_denied"
	ErrorTypeUserNotFound AuthErrorType = "user_not_found"
)

func NewAuthError(errorType AuthErrorType, message string) *AuthError {
	code := http.StatusUnauthorized
	switch errorType {
	case ErrorTypePermission, ErrorTypeRole:
		code = http.StatusForbidden
	case ErrorTypeMissing:
		code = http.StatusUnauthorized
	case ErrorTypeMalformed:
		code = http.StatusBadRequest
	}

	return &AuthError{
		Type:    errorType,
		Message: message,
		Code:    code,
	}
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("auth error [%s]: %s", e.Type, e.Message)
}

func (am *AuthMiddleware) ValidateTokenString(tokenString string) (*CustomClaims, error) {
	if tokenString == "" {
		return nil, NewAuthError(ErrorTypeMissing, "Authorization token is required")
	}

	token, err := am.parseAndValidateToken(tokenString)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, NewAuthError(ErrorTypeExpired, "Token has expired")
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, NewAuthError(ErrorTypeMalformed, "Token is malformed")
		}
		return nil, NewAuthError(ErrorTypeInvalid, "Token is invalid: "+err.Error())
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, NewAuthError(ErrorTypeInvalid, "Invalid token claims")
	}

	if err := am.validateClaims(claims); err != nil {
		if errors.Is(err, ErrTokenExpired) {
			return nil, NewAuthError(ErrorTypeExpired, "Token has expired")
		}
		return nil, NewAuthError(ErrorTypeInvalid, "Token validation failed: "+err.Error())
	}

	return claims, nil
}

func (am *AuthMiddleware) CreateTokenWithDefaults(userID uuid.UUID) (string, error) {
	return am.CreateToken(userID, time.Now().Add(24*time.Hour), map[string]any{})
}

func (am *AuthMiddleware) CreateServiceToken(serviceName string, expiresAt time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub":  serviceName,
		"type": "service",
		"exp":  expiresAt.Unix(),
		"iat":  time.Now().Unix(),
		"nbf":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(am.jwtSecret)
}

func (am *AuthMiddleware) ValidateServiceToken(tokenString string) (serviceName string, err error) {
	if tokenString == "" {
		return "", NewAuthError(ErrorTypeMissing, "Service token is required")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return am.jwtSecret, nil
	})

	if err != nil {
		return "", NewAuthError(ErrorTypeInvalid, "Service token is invalid: "+err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", NewAuthError(ErrorTypeInvalid, "Invalid service token claims")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "service" {
		return "", NewAuthError(ErrorTypeInvalid, "Token is not a service token")
	}

	serviceName, ok = claims["sub"].(string)
	if !ok {
		return "", NewAuthError(ErrorTypeInvalid, "Service name not found in token")
	}

	return serviceName, nil
}

func (am *AuthMiddleware) RefreshToken(tokenString string, newExpiry time.Time) (string, error) {
	claims, err := am.ValidateTokenString(tokenString)
	if err != nil {
		return "", err
	}

	customClaims := map[string]any{
		"permissions": claims.Permissions,
		"roles":       claims.Roles,
	}

	return am.CreateToken(claims.UserID, newExpiry, customClaims)
}

func (am *AuthMiddleware) ExtractUserID(tokenString string) (uuid.UUID, error) {
	claims, err := am.ValidateTokenString(tokenString)
	if err != nil {
		return uuid.Nil, err
	}

	return claims.UserID, nil
}

func (am *AuthMiddleware) ExtractPermissions(tokenString string) ([]string, error) {
	claims, err := am.ValidateTokenString(tokenString)
	if err != nil {
		return nil, err
	}

	return claims.Permissions, nil
}

func (am *AuthMiddleware) ExtractRoles(tokenString string) ([]string, error) {
	claims, err := am.ValidateTokenString(tokenString)
	if err != nil {
		return nil, err
	}

	return claims.Roles, nil
}
