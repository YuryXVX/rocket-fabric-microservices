package user_repository

import (
	"context"
	"iam/internal/model"

	"github.com/google/uuid"
)

func (r *repository) GetUser(ctx context.Context, uuid uuid.UUID) (model.User, error) {
	return model.User{}, ctx.Err()
}
