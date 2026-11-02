package user_repository

import "github.com/jackc/pgx/v5/pgxpool"

type repository struct {
	pool *pgxpool.Pool
}

func New(poll *pgxpool.Pool) *repository {
	return &repository{
		pool: poll,
	}
}
