package session_repository

import (
	"context"
	"fmt"
	"iam/internal/service/input"
)

func (r *repository) Logout(ctx context.Context, in input.LogoutInput) error {
	key := fmt.Sprintf("session:%s", in.SessionUUID.String())

	err := r.client.Del(ctx, key).Err()

	if err != nil {
		return fmt.Errorf("failed to delete session from redis: %w", err)
	}

	return nil
}
