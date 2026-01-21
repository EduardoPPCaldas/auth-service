package token

import (
	"os"
	"testing"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenGenerator_GenerateToken_Success(t *testing.T) {
	// Arrange
	secret := "test-secret-key-for-jwt"
	os.Setenv("JWT_SECRET", secret)
	defer os.Unsetenv("JWT_SECRET")

	generator := NewTokenGenerator()
	user := user.New("test@example.com", nil)

	// Act
	tokenString, err := generator.GenerateToken(user)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Verify token can be parsed
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	require.NoError(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)
	assert.Equal(t, user.ID.String(), claims["sub"])
	assert.NotNil(t, claims["exp"])
}

func TestTokenGenerator_GenerateToken_NoSecret(t *testing.T) {
	// Arrange
	os.Unsetenv("JWT_SECRET")
	generator := NewTokenGenerator()
	user := user.New("test@example.com", nil)

	// Act
	tokenString, err := generator.GenerateToken(user)

	// Assert
	require.Error(t, err)
	assert.Empty(t, tokenString)
}

func TestTokenGenerator_GenerateToken_DifferentUsers(t *testing.T) {
	// Arrange
	secret := "test-secret-key-for-jwt"
	os.Setenv("JWT_SECRET", secret)
	defer os.Unsetenv("JWT_SECRET")

	generator := NewTokenGenerator()
	user1 := user.New("user1@example.com", nil)
	user2 := user.New("user2@example.com", nil)

	// Act
	token1, err1 := generator.GenerateToken(user1)
	token2, err2 := generator.GenerateToken(user2)

	// Assert
	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.NotEqual(t, token1, token2)

	// Verify both tokens contain correct user IDs
	tokenObj1, _ := jwt.Parse(token1, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	tokenObj2, _ := jwt.Parse(token2, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	claims1 := tokenObj1.Claims.(jwt.MapClaims)
	claims2 := tokenObj2.Claims.(jwt.MapClaims)

	assert.Equal(t, user1.ID.String(), claims1["sub"])
	assert.Equal(t, user2.ID.String(), claims2["sub"])
}
