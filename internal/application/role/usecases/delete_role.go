package usecases

import (
	"context"
	"fmt"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/role"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/google/uuid"
)

type DeleteRoleUseCase struct {
	roleRepository role.Repository
	userRepository user.UserRepository
}

func NewDeleteRoleUseCase(roleRepository role.Repository, userRepository user.UserRepository) *DeleteRoleUseCase {
	return &DeleteRoleUseCase{
		roleRepository: roleRepository,
		userRepository: userRepository,
	}
}

type DeleteRoleInput struct {
	AdminUserID uuid.UUID
	RoleID      uuid.UUID
}

func (u *DeleteRoleUseCase) Execute(ctx context.Context, input DeleteRoleInput) error {
	if err := u.verifyAdmin(ctx, input.AdminUserID); err != nil {
		return err
	}

	existingRole, err := u.roleRepository.FindByID(ctx, input.RoleID)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	if existingRole.Name == role.RoleAdmin {
		return fmt.Errorf("cannot delete admin role")
	}

	if existingRole.Name == role.RoleUser {
		return fmt.Errorf("cannot delete default user role")
	}

	if err := u.roleRepository.Delete(ctx, input.RoleID); err != nil {
		return fmt.Errorf("error deleting role: %w", err)
	}

	return nil
}

func (u *DeleteRoleUseCase) verifyAdmin(ctx context.Context, adminUserID uuid.UUID) error {
	adminUser, err := u.userRepository.FindByID(ctx, adminUserID)
	if err != nil {
		return fmt.Errorf("admin user not found: %w", err)
	}

	if adminUser.Role == nil || adminUser.Role.Name != role.RoleAdmin {
		return fmt.Errorf("user does not have admin privileges")
	}

	return nil
}
