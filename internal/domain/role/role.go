package role

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primary_key"`
	Name        string       `json:"name" gorm:"uniqueIndex;not null"`
	Permissions []Permission `json:"permissions" gorm:"foreignKey:RoleID"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type Permission struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Name   string    `json:"name" gorm:"not null"`
	RoleID uuid.UUID `json:"role_id" gorm:"type:uuid;index"`
}

const (
	RoleAdmin     = "admin"
	RoleUser      = "user"
	RoleModerator = "moderator"
)

var (
	AdminPermissions = []string{"*"}

	UserPermissions = []string{
		"users:read:self",
		"users:write:self",
		"posts:read",
		"posts:write:self",
	}

	ModeratorPermissions = []string{
		"posts:read",
		"posts:write",
		"posts:delete",
		"users:read",
	}
)

func New(name string, permissions []string) *Role {
	role := &Role{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, perm := range permissions {
		role.Permissions = append(role.Permissions, Permission{
			ID:     uuid.New(),
			Name:   perm,
			RoleID: role.ID,
		})
	}

	return role
}

func NewAdminRole() *Role {
	return New(RoleAdmin, AdminPermissions)
}

func NewUserRole() *Role {
	return New(RoleUser, UserPermissions)
}

func NewModeratorRole() *Role {
	return New(RoleModerator, ModeratorPermissions)
}

func (r *Role) GetPermissionStrings() []string {
	perms := make([]string, len(r.Permissions))
	for i, p := range r.Permissions {
		perms[i] = p.Name
	}
	return perms
}

func (r *Role) HasPermission(permission string) bool {
	if r.Name == RoleAdmin {
		return true
	}

	for _, p := range r.Permissions {
		if p.Name == permission || p.Name == "*" {
			return true
		}
	}
	return false
}
