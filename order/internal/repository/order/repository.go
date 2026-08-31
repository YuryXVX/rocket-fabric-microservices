package order

import (
	"github.com/jackc/pgx/v5/pgxpool"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
)

type repository struct {
	getter *trmpgx.CtxGetter
	pool   *pgxpool.Pool
	tx     TxManager
}

func New(pool *pgxpool.Pool, tx TxManager) *repository {
	return &repository{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
		tx:     tx,
	}
}
