package order

import (
	"context"
	"errors"
	errs "order/internal/errors"
	"order/internal/model"
	"order/internal/repository/converter"
	"order/internal/repository/record"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *repository) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	var rec record.OrderRecord

	err := r.getter.DefaultTrOrDB(ctx, r.pool).
		QueryRow(ctx, "SELECT * FROM orders WHERE uuid = ($1)", orderUUID).
		Scan(&rec.OrderUUID, &rec.Status, &rec.TransactionUUID, &rec.PaymentMethod, &rec.CreatedAt, &rec.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, errs.ErrOrderNotFound
		}
		return model.Order{}, err
	}

	return converter.RecordOrderToModel(rec), nil
}
