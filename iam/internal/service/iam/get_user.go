package iam_service

import (
	"context"
	"iam/internal/model"

	"github.com/google/uuid"
)

func (s *service) GetUser(ctx context.Context, userID uuid.UUID) (model.User, error) {
	// Implementation here
	return model.User{}, nil
}
