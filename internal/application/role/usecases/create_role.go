package usecases

import (
	"context"
	"fmt"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/role"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateRoleUseCase struct {
	roleRepository role.Repository
	userRepository user.UserRepository
}

func NewCreateRoleUseCase(roleRepository role.Repository, userRepository user.UserRepository) *CreateRoleUseCase {
	return &CreateRoleUseCase{
		roleRepository: roleRepository,
		userRepository: userRepository,
	}
}

type CreateRoleInput struct {
	AdminUserID uuid.UUID
	Name        string
	Permissions []string
}

func (u *CreateRoleUseCase) Execute(ctx context.Context, input CreateRoleInput) (*role.Role, error) {
	if err := u.verifyAdmin(ctx, input.AdminUserID); err != nil {
		return nil, err
	}

	existingRole, err := u.roleRepository.FindByName(ctx, input.Name)
	if err != nil && !isNotFound(err) {
		return nil, fmt.Errorf("error checking existing role: %w", err)
	}
	if existingRole != nil {
		return nil, fmt.Errorf("role with name '%s' already exists", input.Name)
	}

	newRole := role.New(input.Name, input.Permissions)
	if err := u.roleRepository.Create(ctx, newRole); err != nil {
		return nil, fmt.Errorf("error creating role: %w", err)
	}

	return newRole, nil
}

func (u *CreateRoleUseCase) verifyAdmin(ctx context.Context, adminUserID uuid.UUID) error {
	adminUser, err := u.userRepository.FindByID(ctx, adminUserID)
	if err != nil {
		return fmt.Errorf("admin user not found: %w", err)
	}

	if adminUser.Role == nil || adminUser.Role.Name != role.RoleAdmin {
		return fmt.Errorf("user does not have admin privileges")
	}

	return nil
}

func isNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
