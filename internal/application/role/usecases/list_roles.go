package usecases

import (
	"context"
	"fmt"

	"github.com/EduardoPPCaldas/auth-service/internal/domain/role"
)

type ListRolesUseCase struct {
	roleRepository role.Repository
}

func NewListRolesUseCase(roleRepository role.Repository) *ListRolesUseCase {
	return &ListRolesUseCase{
		roleRepository: roleRepository,
	}
}

func (u *ListRolesUseCase) Execute(ctx context.Context) ([]role.Role, error) {
	roles, err := u.roleRepository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing roles: %w", err)
	}
	return roles, nil
}
