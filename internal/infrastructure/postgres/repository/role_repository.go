package repository

import (
	"context"
	"errors"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/role"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleRepository implements role.Repository interface
type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// Create creates a new role using GORM generics
func (r *RoleRepository) Create(ctx context.Context, ro *role.Role) error {
	return gorm.G[role.Role](r.db).Create(ctx, ro)
}

// FindByID finds a role by its ID with permissions preloaded
func (r *RoleRepository) FindByID(ctx context.Context, id uuid.UUID) (*role.Role, error) {
	var ro role.Role
	result := r.db.WithContext(ctx).Preload("Permissions").Where("id = ?", id).First(&ro)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ro, nil
}

// FindByName finds a role by its name with permissions preloaded
func (r *RoleRepository) FindByName(ctx context.Context, name string) (*role.Role, error) {
	ro, err := gorm.G[role.Role](r.db.WithContext(ctx)).Preload("Permissions", nil).Where("name = ?", name).First(ctx)
	return &ro, err
}

// FindOrCreateDefault returns the default user role if RBAC is enabled
func (r *RoleRepository) FindOrCreateDefault(ctx context.Context) (*role.Role, error) {
	existingRole, err := r.FindByName(ctx, role.RoleUser)
	if err == nil {
		return existingRole, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// RBAC not enabled, return nil
	return nil, nil
}

// IsRBACEnabled checks if any roles exist in the database
func (r *RoleRepository) IsRBACEnabled(ctx context.Context) bool {
	var count int64
	r.db.WithContext(ctx).Model(&role.Role{}).Count(&count)
	return count > 0
}

// SeedRoles creates default roles if they don't exist
func (r *RoleRepository) SeedRoles(ctx context.Context) error {
	roles := []*role.Role{
		role.NewAdminRole(),
		role.NewUserRole(),
		role.NewModeratorRole(),
	}

	for _, ro := range roles {
		var existing role.Role
		r.db.WithContext(ctx).Where("name = ?", ro.Name).FirstOrCreate(&existing)
	}

	return nil
}
