package middleware

import (
	"context"
	"errors"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

func RequireAdmin(ctx context.Context) error {
	role, ok := GetRoleFromContext(ctx)
	if !ok {
		return errors.New("role not found in context")
	}

	if role != RoleAdmin {
		return errors.New("admin role required")
	}

	return nil
}

func RequireUser(ctx context.Context) error {
	role, ok := GetRoleFromContext(ctx)
	if !ok {
		return errors.New("role not found in context")
	}

	if role != RoleUser && role != RoleAdmin {
		return errors.New("user role required")
	}

	return nil
}

func HasRole(ctx context.Context, allowedRoles ...string) bool {
	role, ok := GetRoleFromContext(ctx)
	if !ok {
		return false
	}

	for _, allowedRole := range allowedRoles {
		if role == allowedRole {
			return true
		}
	}

	return false
}
