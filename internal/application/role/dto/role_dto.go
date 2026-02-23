package dto

import "github.com/EduardoPPCaldas/auth-service/internal/domain/role"

type CreateRoleRequest struct {
	Name        string   `json:"name" validate:"required,min=2,max=50"`
	Permissions []string `json:"permissions" validate:"required,min=1,dive,required"`
}

type UpdateRoleRequest struct {
	Name        string   `json:"name" validate:"omitempty,min=2,max=50"`
	Permissions []string `json:"permissions" validate:"omitempty,min=1,dive,required"`
}

type RoleResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Permissions []PermissionResponse `json:"permissions"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
}

type PermissionResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ToRoleResponse(r *role.Role) RoleResponse {
	perms := make([]PermissionResponse, len(r.Permissions))
	for i, p := range r.Permissions {
		perms[i] = PermissionResponse{
			ID:   p.ID.String(),
			Name: p.Name,
		}
	}
	return RoleResponse{
		ID:          r.ID.String(),
		Name:        r.Name,
		Permissions: perms,
		CreatedAt:   r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   r.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

type AssignRoleRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	RoleID string `json:"role_id" validate:"required,uuid"`
}
