package session_repository

import (
	"context"
	"iam/internal/service/input"

	"github.com/google/uuid"
)

func (r *repository) Login(ctx context.Context, in input.LoginInput) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}
