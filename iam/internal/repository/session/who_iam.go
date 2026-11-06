package session_repository

import (
	"context"
	"errors"
	"fmt"
	"iam/internal/model"
	"iam/internal/repository/converter"
	"iam/internal/repository/redis_view"

	"github.com/google/uuid"
)

var ErrSessionNotFound = errors.New("session not found in cache")

func (r *repository) WhoIam(ctx context.Context, in uuid.UUID) (model.Session, error) {
	key := fmt.Sprintf("session:%s", in.String())

	var view redis_view.SessionView

	err := r.client.HGetAll(ctx, key).Scan(&view)

	if err != nil {
		return model.Session{}, err
	}

	if view.UserUUID == uuid.Nil {
		return model.Session{}, ErrSessionNotFound
	}

	return converter.SessionRedisViewToModel(view), nil
}
