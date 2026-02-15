package repository

import (
	"context"

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
func (r *RoleRepository) FindByID(id uuid.UUID) (*role.Role, error) {
	var ro role.Role
	result := r.db.Preload("Permissions").Where("id = ?", id).First(&ro)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ro, nil
}

// FindByName finds a role by its name with permissions preloaded
func (r *RoleRepository) FindByName(name string) (*role.Role, error) {
	var ro role.Role
	result := r.db.Preload("Permissions").Where("name = ?", name).First(&ro)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ro, nil
}

// FindOrCreateDefault returns the default user role if RBAC is enabled
func (r *RoleRepository) FindOrCreateDefault() (*role.Role, error) {
	existingRole, err := r.FindByName(role.RoleUser)
	if err == nil {
		return existingRole, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// RBAC not enabled, return nil
	return nil, nil
}

// IsRBACEnabled checks if any roles exist in the database
func (r *RoleRepository) IsRBACEnabled() bool {
	var count int64
	r.db.Model(&role.Role{}).Count(&count)
	return count > 0
}

// SeedRoles creates default roles if they don't exist
func (r *RoleRepository) SeedRoles() error {
	ctx := context.Background()
	roles := []*role.Role{
		role.NewAdminRole(),
		role.NewUserRole(),
		role.NewModeratorRole(),
	}

	for _, ro := range roles {
		var existing role.Role
		result := r.db.Where("name = ?", ro.Name).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := r.Create(ctx, ro); err != nil {
				return err
			}
		}
	}

	return nil
}
