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

// Update updates an existing role
func (r *RoleRepository) Update(ctx context.Context, ro *role.Role) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update role basic info
		if err := tx.Model(ro).Updates(map[string]any{
			"name":       ro.Name,
			"updated_at": ro.UpdatedAt,
		}).Error; err != nil {
			return err
		}
		// Delete old permissions and create new ones
		if err := tx.Where("role_id = ?", ro.ID).Delete(&role.Permission{}).Error; err != nil {
			return err
		}
		if len(ro.Permissions) > 0 {
			if err := tx.Create(ro.Permissions).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Delete deletes a role by ID
func (r *RoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete permissions first
		if err := tx.Where("role_id = ?", id).Delete(&role.Permission{}).Error; err != nil {
			return err
		}
		// Delete role
		return tx.Delete(&role.Role{}, "id = ?", id).Error
	})
}

// List returns all roles with their permissions
func (r *RoleRepository) List(ctx context.Context) ([]role.Role, error) {
	var roles []role.Role
	result := r.db.WithContext(ctx).Preload("Permissions").Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}
