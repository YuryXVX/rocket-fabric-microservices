package user_repository

import (
	"context"
	"errors"
	"fmt"
	"iam/internal/model"

	"github.com/jackc/pgx/v5"
)

func (r *repository) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	var u model.User

	query := "SELECT uuid, login, password_hash, created_at, updated_at FROM users WHERE login = $1"

	err := r.pool.
		QueryRow(ctx, query, login).
		Scan(
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
