package user_repository

import (
	"context"
	"errors"
	"iam/internal/model"
	"iam/internal/repository/converter"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrUserAlreadyExists = errors.New("такой пользователь уже сущевует")

func (r *repository) Register(ctx context.Context, model model.User) (uuid.UUID, error) {
	rec := converter.UserModelToRecord(model)

	_, err := r.pool.Exec(ctx,
		"INSERT INTO users (uuid, login, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		rec.UUID,
		rec.Login,
		rec.PasswordHash,
		rec.CreatedAt,
		rec.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return uuid.UUID{}, ErrUserAlreadyExists
		}
		return uuid.UUID{}, err
	}

	return model.UUID, nil
}
