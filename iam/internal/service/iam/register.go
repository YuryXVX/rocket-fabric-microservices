package iam_service

import (
	"context"
	"iam/internal/service/input"

	"github.com/google/uuid"
)

func (s *service) Register(ctx context.Context, in input.RegisterInput) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}
