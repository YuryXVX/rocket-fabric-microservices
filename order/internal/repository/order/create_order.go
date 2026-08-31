package order

import (
	"context"
	"fmt"
	"order/internal/model"
	"order/internal/repository/converter"
)

func (r *repository) createOrder(ctx context.Context, order model.Order) error {
	rec := converter.OrderModelToRecord(&order)

	_, err := r.getter.
		DefaultTrOrDB(ctx, r.pool).
		Exec(ctx,
			"INSERT INTO orders (uuid, status, transaction_uuid, payment_method, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
			rec.OrderUUID,
			rec.Status,
			rec.TransactionUUID,
			rec.PaymentMethod,
			rec.CreatedAt,
			rec.UpdatedAt,
		)

	if err != nil {
		return fmt.Errorf("создать заказ: %w", err)
	}

	return nil
}
