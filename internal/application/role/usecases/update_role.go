package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/role"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/google/uuid"
)

type UpdateRoleUseCase struct {
	roleRepository role.Repository
	userRepository user.UserRepository
}

func NewUpdateRoleUseCase(roleRepository role.Repository, userRepository user.UserRepository) *UpdateRoleUseCase {
	return &UpdateRoleUseCase{
		roleRepository: roleRepository,
		userRepository: userRepository,
	}
}

type UpdateRoleInput struct {
	AdminUserID uuid.UUID
	RoleID      uuid.UUID
	Name        *string
	Permissions []string
}

func (u *UpdateRoleUseCase) Execute(ctx context.Context, input UpdateRoleInput) (*role.Role, error) {
	if err := u.verifyAdmin(ctx, input.AdminUserID); err != nil {
		return nil, err
	}

	existingRole, err := u.roleRepository.FindByID(ctx, input.RoleID)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}

	if existingRole.Name == role.RoleAdmin {
		return nil, fmt.Errorf("cannot modify admin role")
	}

	if input.Name != nil && *input.Name != existingRole.Name {
		roleWithName, err := u.roleRepository.FindByName(ctx, *input.Name)
		if err != nil && !isNotFound(err) {
			return nil, fmt.Errorf("error checking existing role: %w", err)
		}
		if roleWithName != nil {
			return nil, fmt.Errorf("role with name '%s' already exists", *input.Name)
		}
		existingRole.Name = *input.Name
	}

	if input.Permissions != nil {
		existingRole.Permissions = nil
		for _, perm := range input.Permissions {
			existingRole.Permissions = append(existingRole.Permissions, role.Permission{
				ID:     uuid.New(),
				Name:   perm,
				RoleID: existingRole.ID,
			})
		}
	}

	existingRole.UpdatedAt = time.Now()

	if err := u.roleRepository.Update(ctx, existingRole); err != nil {
		return nil, fmt.Errorf("error updating role: %w", err)
	}

	return existingRole, nil
}

func (u *UpdateRoleUseCase) verifyAdmin(ctx context.Context, adminUserID uuid.UUID) error {
	adminUser, err := u.userRepository.FindByID(ctx, adminUserID)
	if err != nil {
		return fmt.Errorf("admin user not found: %w", err)
	}

	if adminUser.Role == nil || adminUser.Role.Name != role.RoleAdmin {
		return fmt.Errorf("user does not have admin privileges")
	}

	return nil
}
