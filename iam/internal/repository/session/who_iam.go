package session_repository

import (
	"context"
	"errors"
	"fmt"
	"iam/internal/model"
	"iam/internal/repository/converter"
	"iam/internal/repository/redis_view"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var ErrSessionNotFound = errors.New("session not found in cache")

func (r *repository) WhoIam(ctx context.Context, in uuid.UUID) (model.Session, error) {
	key := fmt.Sprintf("session:%s", in.String())

	var view redis_view.SessionView

	err := r.client.HGetAll(ctx, key).Scan(&view)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.Session{}, ErrSessionNotFound
		}

		return model.Session{}, err
	}

	view.UUID = in

	return converter.SessionRedisViewToModel(view), nil
}
