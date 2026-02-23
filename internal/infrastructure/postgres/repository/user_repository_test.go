package repository

import (
	"context"
	"testing"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&user.User{})
	require.NoError(t, err)

	return db
}

func TestUserRepository_Create(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	testUser := user.New("test@example.com", nil)
	ctx := context.Background()

	// Act
	err := repo.Create(ctx, testUser)

	// Assert
	require.NoError(t, err)

	var foundUser user.User
	err = db.First(&foundUser, "id = ?", testUser.ID).Error
	require.NoError(t, err)
	assert.Equal(t, testUser.Email, foundUser.Email)
	assert.Equal(t, testUser.ID, foundUser.ID)
}

func TestUserRepository_FindByEmail_Exists(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	testUser := user.New("test@example.com", nil)
	ctx := context.Background()

	err := db.Create(testUser).Error
	require.NoError(t, err)

	// Act
	foundUser, err := repo.FindByEmail(ctx, "test@example.com")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, testUser.Email, foundUser.Email)
	assert.Equal(t, testUser.ID, foundUser.ID)
}

func TestUserRepository_FindByEmail_NotExists(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Act
	foundUser, err := repo.FindByEmail(ctx, "nonexistent@example.com")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, foundUser)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUserRepository_FindByEmail_CaseSensitive(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	testUser := user.New("Test@Example.com", nil)
	ctx := context.Background()

	err := db.Create(testUser).Error
	require.NoError(t, err)

	// Act
	foundUser1, err1 := repo.FindByEmail(ctx, "Test@Example.com")
	foundUser2, err2 := repo.FindByEmail(ctx, "test@example.com")

	// Assert
	require.NoError(t, err1)
	assert.NotNil(t, foundUser1)
	assert.Equal(t, testUser.Email, foundUser1.Email)

	// Email comparison is case-sensitive in SQL
	assert.Error(t, err2)
	assert.Nil(t, foundUser2)
}
