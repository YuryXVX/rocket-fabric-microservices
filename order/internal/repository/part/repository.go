package part

import (
	"github.com/jackc/pgx/v5/pgxpool"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
)

type repository struct {
	getter *trmpgx.CtxGetter
	pool   *pgxpool.Pool
}

func New(poll *pgxpool.Pool) *repository {
	return &repository{
		pool:   poll,
		getter: trmpgx.DefaultCtxGetter,
	}
}
