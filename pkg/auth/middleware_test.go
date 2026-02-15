package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAuthMiddleware(t *testing.T) {
	secret := "test-secret-key"
	authMiddleware, err := NewAuthMiddleware(WithJWTSecret(secret))
	if err != nil {
		t.Fatalf("Failed to create auth middleware: %v", err)
	}

	userID := uuid.New()

	// Test token creation
	token, err := authMiddleware.CreateTokenWithDefaults(userID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	if token == "" {
		t.Fatal("Token is empty")
	}

	// Test token validation
	claims, err := authMiddleware.ValidateTokenString(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, claims.UserID)
	}

	// Test extracting user ID
	extractedUserID, err := authMiddleware.ExtractUserID(token)
	if err != nil {
		t.Fatalf("Failed to extract user ID: %v", err)
	}

	if extractedUserID != userID {
		t.Errorf("Expected extracted user ID %s, got %s", userID, extractedUserID)
	}
}

func TestServiceToken(t *testing.T) {
	secret := "test-secret-key"
	authMiddleware, err := NewAuthMiddleware(WithJWTSecret(secret))
	if err != nil {
		t.Fatalf("Failed to create auth middleware: %v", err)
	}

	serviceName := "test-service"
	expiresAt := time.Now().Add(1 * time.Hour)

	// Create service token
	token, err := authMiddleware.CreateServiceToken(serviceName, expiresAt)
	if err != nil {
		t.Fatalf("Failed to create service token: %v", err)
	}

	// Validate service token
	extractedServiceName, err := authMiddleware.ValidateServiceToken(token)
	if err != nil {
		t.Fatalf("Failed to validate service token: %v", err)
	}

	if extractedServiceName != serviceName {
		t.Errorf("Expected service name %s, got %s", serviceName, extractedServiceName)
	}
}

func TestTokenWithCustomClaims(t *testing.T) {
	secret := "test-secret-key"
	authMiddleware, err := NewAuthMiddleware(WithJWTSecret(secret))
	if err != nil {
		t.Fatalf("Failed to create auth middleware: %v", err)
	}

	userID := uuid.New()
	customClaims := map[string]any{
		"permissions": []string{"read:users", "write:users"},
		"roles":       []string{"admin", "moderator"},
	}

	token, err := authMiddleware.CreateToken(userID, time.Now().Add(1*time.Hour), customClaims)
	if err != nil {
		t.Fatalf("Failed to create token with custom claims: %v", err)
	}

	// Test extracting permissions
	permissions, err := authMiddleware.ExtractPermissions(token)
	if err != nil {
		t.Fatalf("Failed to extract permissions: %v", err)
	}

	expectedPerms := []string{"read:users", "write:users"}
	if len(permissions) != len(expectedPerms) {
		t.Errorf("Expected %d permissions, got %d", len(expectedPerms), len(permissions))
	}

	for i, perm := range expectedPerms {
		if i >= len(permissions) || permissions[i] != perm {
			t.Errorf("Expected permission %s, got %v", perm, permissions)
		}
	}

	// Test extracting roles
	roles, err := authMiddleware.ExtractRoles(token)
	if err != nil {
		t.Fatalf("Failed to extract roles: %v", err)
	}

	expectedRoles := []string{"admin", "moderator"}
	if len(roles) != len(expectedRoles) {
		t.Errorf("Expected %d roles, got %d", len(expectedRoles), len(roles))
	}

	for i, role := range expectedRoles {
		if i >= len(roles) || roles[i] != role {
			t.Errorf("Expected role %s, got %v", role, roles)
		}
	}
}

func TestTokenRefresh(t *testing.T) {
	secret := "test-secret-key"
	authMiddleware, err := NewAuthMiddleware(WithJWTSecret(secret))
	if err != nil {
		t.Fatalf("Failed to create auth middleware: %v", err)
	}

	userID := uuid.New()
	originalToken, err := authMiddleware.CreateTokenWithDefaults(userID)
	if err != nil {
		t.Fatalf("Failed to create original token: %v", err)
	}

	// Refresh token with new expiry
	newExpiry := time.Now().Add(2 * time.Hour)
	refreshedToken, err := authMiddleware.RefreshToken(originalToken, newExpiry)
	if err != nil {
		t.Fatalf("Failed to refresh token: %v", err)
	}

	if refreshedToken == originalToken {
		t.Error("Refreshed token should be different from original token")
	}

	// Validate refreshed token
	claims, err := authMiddleware.ValidateTokenString(refreshedToken)
	if err != nil {
		t.Fatalf("Failed to validate refreshed token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected user ID %s in refreshed token, got %s", userID, claims.UserID)
	}
}

func TestInvalidTokens(t *testing.T) {
	secret := "test-secret-key"
	authMiddleware, err := NewAuthMiddleware(WithJWTSecret(secret))
	if err != nil {
		t.Fatalf("Failed to create auth middleware: %v", err)
	}

	// Test empty token
	_, err = authMiddleware.ValidateTokenString("")
	if err == nil {
		t.Error("Expected error for empty token")
	}

	// Test malformed token
	_, err = authMiddleware.ValidateTokenString("invalid.token.here")
	if err == nil {
		t.Error("Expected error for malformed token")
	}

	// Test token with wrong secret
	wrongAuthMiddleware, _ := NewAuthMiddleware(WithJWTSecret("wrong-secret"))
	validToken, _ := authMiddleware.CreateTokenWithDefaults(uuid.New())

	_, err = wrongAuthMiddleware.ValidateTokenString(validToken)
	if err == nil {
		t.Error("Expected error for token with wrong secret")
	}
}
