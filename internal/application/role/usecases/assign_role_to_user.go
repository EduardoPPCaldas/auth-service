package usecases

import (
	"context"
	"fmt"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/role"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/google/uuid"
)

type AssignRoleToUserUseCase struct {
	roleRepository role.Repository
	userRepository user.UserRepository
}

func NewAssignRoleToUserUseCase(roleRepository role.Repository, userRepository user.UserRepository) *AssignRoleToUserUseCase {
	return &AssignRoleToUserUseCase{
		roleRepository: roleRepository,
		userRepository: userRepository,
	}
}

type AssignRoleToUserInput struct {
	AdminUserID uuid.UUID
	UserID      uuid.UUID
	RoleID      uuid.UUID
}

func (u *AssignRoleToUserUseCase) Execute(ctx context.Context, input AssignRoleToUserInput) error {
	if err := u.verifyAdmin(ctx, input.AdminUserID); err != nil {
		return err
	}

	targetUser, err := u.userRepository.FindByID(ctx, input.UserID)
	if err != nil {
		return fmt.Errorf("target user not found: %w", err)
	}

	ro, err := u.roleRepository.FindByID(ctx, input.RoleID)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	if targetUser.Role != nil && targetUser.Role.ID == ro.ID {
		return fmt.Errorf("user already has this role assigned")
	}

	if err := u.userRepository.UpdateRole(ctx, input.UserID, &ro.ID); err != nil {
		return fmt.Errorf("error assigning role to user: %w", err)
	}

	return nil
}

func (u *AssignRoleToUserUseCase) verifyAdmin(ctx context.Context, adminUserID uuid.UUID) error {
	adminUser, err := u.userRepository.FindByID(ctx, adminUserID)
	if err != nil {
		return fmt.Errorf("admin user not found: %w", err)
	}

	if adminUser.Role == nil || adminUser.Role.Name != role.RoleAdmin {
		return fmt.Errorf("user does not have admin privileges")
	}

	return nil
}
