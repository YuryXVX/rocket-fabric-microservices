package user_repository

import (
	"context"
	"errors"
	"fmt"
	"iam/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrUserNotFound = errors.New("user not found in database")

func (r *repository) GetUser(ctx context.Context, id uuid.UUID) (model.User, error) {
	var u model.User

	query := "SELECT uuid, login, password_hash, created_at, updated_at FROM users WHERE uuid = $1"

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.UUID,
		&u.Login,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}

		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return u, nil
}
