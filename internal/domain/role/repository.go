package role

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	FindByID(id uuid.UUID) (*Role, error)
	FindByName(name string) (*Role, error)
	Create(ctx context.Context, role *Role) error
	FindOrCreateDefault() (*Role, error)
	IsRBACEnabled() bool
}
