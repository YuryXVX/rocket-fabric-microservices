package session_repository

import (
	"context"
	"iam/internal/model"

	"github.com/google/uuid"
)

func (r *repository) WhoIam(ctx context.Context, in uuid.UUID) (model.Session, error) {
	return model.Session{}, nil
}
