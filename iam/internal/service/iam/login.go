package iam_service

import (
	"context"
	"iam/internal/service/input"

	"github.com/google/uuid"
)

func (s *service) Login(context.Context, input.LoginInput) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}
