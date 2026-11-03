package session_repository

import (
	"context"
	"iam/internal/service/input"
)

func (r *repository) Logout(ctx context.Context, in input.LogoutInput) error {
	return nil
}
