package repository

import (
	"context"
	"testing"
	"time"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/role"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupIntegrationDB(t *testing.T) (*gorm.DB, func()) {
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

	err = db.AutoMigrate(&role.Role{}, &role.Permission{}, &user.User{})
	require.NoError(t, err)

	cleanup := func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}

	return db, cleanup
}

func TestUserRepository_Integration_Create(t *testing.T) {
	// Arrange
	ctx := context.Background()
	db, cleanup := setupIntegrationDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	roleRepo := NewRoleRepository(db)

	roleRepo.Create(ctx, &role.Role{
		Name: role.RoleUser,
	})

	defaultRole, err := roleRepo.FindByName(ctx, role.RoleUser)
	require.NoError(t, err)

	testUser := user.New("integration@example.com", nil)
	testUser.RoleID = &defaultRole.ID

	// Act
	err = repo.Create(ctx, testUser)

	// Assert
	require.NoError(t, err)

	var foundUser user.User
	err = db.Preload("Role").Preload("Role.Permissions").First(&foundUser, "id = ?", testUser.ID).Error
	require.NoError(t, err)
	assert.Equal(t, testUser.Email, foundUser.Email)
	assert.Equal(t, testUser.ID, foundUser.ID)
	assert.Equal(t, defaultRole.ID, *foundUser.RoleID)
}

func TestUserRepository_Integration_FindByEmail(t *testing.T) {
	// Arrange
	ctx := context.Background()
	db, cleanup := setupIntegrationDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	roleRepo := NewRoleRepository(db)

	roleRepo.Create(ctx, &role.Role{
		Name: role.RoleUser,
	})

	defaultRole, err := roleRepo.FindByName(ctx, role.RoleUser)
	require.NoError(t, err)

	testUser := user.New("findtest@example.com", nil)
	testUser.RoleID = &defaultRole.ID

	err = db.Create(testUser).Error
	require.NoError(t, err)

	// Act
	foundUser, err := repo.FindByEmail(ctx, "findtest@example.com")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, testUser.Email, foundUser.Email)
	assert.Equal(t, testUser.ID, foundUser.ID)
	assert.Equal(t, defaultRole.ID, *foundUser.RoleID)
}

func TestUserRepository_Integration_CreateAndFind(t *testing.T) {
	// Arrange
	ctx := context.Background()
	db, cleanup := setupIntegrationDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	roleRepo := NewRoleRepository(db)

	roleRepo.Create(ctx, &role.Role{
		Name: role.RoleUser,
	})

	defaultRole, err := roleRepo.FindByName(ctx, role.RoleUser)
	require.NoError(t, err)

	testUser := user.New("createfind@example.com", nil)
	testUser.RoleID = &defaultRole.ID

	// Act - Create
	err = repo.Create(ctx, testUser)
	require.NoError(t, err)

	// Act - Find
	foundUser, err := repo.FindByEmail(ctx, "createfind@example.com")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, testUser.Email, foundUser.Email)
	assert.Equal(t, testUser.ID, foundUser.ID)
	assert.Equal(t, defaultRole.ID, *foundUser.RoleID)
}

func TestUserRepository_Integration_DuplicateEmail(t *testing.T) {
	// Arrange
	ctx := context.Background()
	db, cleanup := setupIntegrationDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	roleRepo := NewRoleRepository(db)

	roleRepo.Create(ctx, &role.Role{
		Name: role.RoleUser,
	})

	defaultRole, err := roleRepo.FindByName(ctx, role.RoleUser)
	require.NoError(t, err)

	testUser1 := user.New("duplicate@example.com", nil)
	testUser1.RoleID = &defaultRole.ID
	testUser2 := user.New("duplicate@example.com", nil)
	testUser2.RoleID = &defaultRole.ID

	// Act
	err1 := repo.Create(ctx, testUser1)
	err2 := repo.Create(ctx, testUser2)

	// Assert
	require.NoError(t, err1)
	assert.Error(t, err2) // Should fail due to unique constraint (if enabled)
}
