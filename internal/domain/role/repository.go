package role

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Role, error)
	FindByName(ctx context.Context, name string) (*Role, error)
	Create(ctx context.Context, role *Role) error
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]Role, error)
	FindOrCreateDefault(ctx context.Context) (*Role, error)
	IsRBACEnabled(ctx context.Context) bool
}
