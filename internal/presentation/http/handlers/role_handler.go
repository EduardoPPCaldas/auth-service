package handlers

import (
	"context"
	"net/http"

	"github.com/EduardoPPCaldas/auth-service/internal/application/role/dto"
	roleusecases "github.com/EduardoPPCaldas/auth-service/internal/application/role/usecases"
	"github.com/EduardoPPCaldas/auth-service/internal/domain/role"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
	createRoleUseCase       CreateRoleUseCase
	updateRoleUseCase       UpdateRoleUseCase
	deleteRoleUseCase       DeleteRoleUseCase
	listRolesUseCase        ListRolesUseCase
	getRoleUseCase          GetRoleUseCase
	assignRoleToUserUseCase AssignRoleToUserUseCase
}

type CreateRoleUseCase interface {
	Execute(ctx context.Context, input roleusecases.CreateRoleInput) (*role.Role, error)
}

type UpdateRoleUseCase interface {
	Execute(ctx context.Context, input roleusecases.UpdateRoleInput) (*role.Role, error)
}

type DeleteRoleUseCase interface {
	Execute(ctx context.Context, input roleusecases.DeleteRoleInput) error
}

type ListRolesUseCase interface {
	Execute(ctx context.Context) ([]role.Role, error)
}

type GetRoleUseCase interface {
	Execute(ctx context.Context, roleID uuid.UUID) (*role.Role, error)
}

type AssignRoleToUserUseCase interface {
	Execute(ctx context.Context, input roleusecases.AssignRoleToUserInput) error
}

func NewRoleHandler(
	createRoleUseCase CreateRoleUseCase,
	updateRoleUseCase UpdateRoleUseCase,
	deleteRoleUseCase DeleteRoleUseCase,
	listRolesUseCase ListRolesUseCase,
	getRoleUseCase GetRoleUseCase,
	assignRoleToUserUseCase AssignRoleToUserUseCase,
) *RoleHandler {
	return &RoleHandler{
		createRoleUseCase:       createRoleUseCase,
		updateRoleUseCase:       updateRoleUseCase,
		deleteRoleUseCase:       deleteRoleUseCase,
		listRolesUseCase:        listRolesUseCase,
		getRoleUseCase:          getRoleUseCase,
		assignRoleToUserUseCase: assignRoleToUserUseCase,
	}
}

// CreateRole handles role creation
// POST /api/v1/admin/roles
func (h *RoleHandler) CreateRole(c echo.Context) error {
	var req dto.CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	adminUserIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not authenticated"})
	}

	adminUserID, err := uuid.Parse(adminUserIDStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user ID format"})
	}

	input := roleusecases.CreateRoleInput{
		AdminUserID: adminUserID,
		Name:        req.Name,
		Permissions: req.Permissions,
	}

	ro, err := h.createRoleUseCase.Execute(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.ToRoleResponse(ro))
}

// UpdateRole handles role updates
// PUT /api/v1/admin/roles/:id
func (h *RoleHandler) UpdateRole(c echo.Context) error {
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid role ID"})
	}

	var req dto.UpdateRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	adminUserIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not authenticated"})
	}

	adminUserID, err := uuid.Parse(adminUserIDStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user ID format"})
	}

	input := roleusecases.UpdateRoleInput{
		AdminUserID: adminUserID,
		RoleID:      roleID,
		Permissions: req.Permissions,
	}

	if req.Name != "" {
		input.Name = &req.Name
	}

	ro, err := h.updateRoleUseCase.Execute(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ToRoleResponse(ro))
}

// DeleteRole handles role deletion
// DELETE /api/v1/admin/roles/:id
func (h *RoleHandler) DeleteRole(c echo.Context) error {
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid role ID"})
	}

	adminUserIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not authenticated"})
	}

	adminUserID, err := uuid.Parse(adminUserIDStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user ID format"})
	}

	input := roleusecases.DeleteRoleInput{
		AdminUserID: adminUserID,
		RoleID:      roleID,
	}

	if err := h.deleteRoleUseCase.Execute(c.Request().Context(), input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "role deleted successfully"})
}

// ListRoles handles listing all roles
// GET /api/v1/admin/roles
func (h *RoleHandler) ListRoles(c echo.Context) error {
	roles, err := h.listRolesUseCase.Execute(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := make([]dto.RoleResponse, len(roles))
	for i, ro := range roles {
		response[i] = dto.ToRoleResponse(&ro)
	}

	return c.JSON(http.StatusOK, response)
}

// GetRole handles getting a specific role
// GET /api/v1/admin/roles/:id
func (h *RoleHandler) GetRole(c echo.Context) error {
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid role ID"})
	}

	ro, err := h.getRoleUseCase.Execute(c.Request().Context(), roleID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ToRoleResponse(ro))
}

// AssignRoleToUser handles assigning a role to a user
// POST /api/v1/admin/roles/assign
func (h *RoleHandler) AssignRoleToUser(c echo.Context) error {
	var req dto.AssignRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	adminUserIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not authenticated"})
	}

	adminUserID, err := uuid.Parse(adminUserIDStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user ID format"})
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid role ID"})
	}

	input := roleusecases.AssignRoleToUserInput{
		AdminUserID: adminUserID,
		UserID:      userID,
		RoleID:      roleID,
	}

	if err := h.assignRoleToUserUseCase.Execute(c.Request().Context(), input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "role assigned successfully"})
}
