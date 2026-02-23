package usecases

import (
	"context"
	"fmt"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/role"
	"github.com/google/uuid"
)

type GetRoleUseCase struct {
	roleRepository role.Repository
}

func NewGetRoleUseCase(roleRepository role.Repository) *GetRoleUseCase {
	return &GetRoleUseCase{
		roleRepository: roleRepository,
	}
}

func (u *GetRoleUseCase) Execute(ctx context.Context, roleID uuid.UUID) (*role.Role, error) {
	ro, err := u.roleRepository.FindByID(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}
	return ro, nil
}
