package order

import (
	"context"
	"fmt"
	"order/internal/model"
)

func (r *repository) Update(ctx context.Context, order model.Order) error {
	const sql = `
    UPDATE orders 
    SET status = $2, 
        transaction_uuid = $3, 
        payment_method = $4, 
        created_at = $5, 
        updated_at = $6 
    WHERE uuid = $1`

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).
		Exec(ctx, sql,
			order.UUID,
			order.Status,
			order.TransactionUUID,
			order.PaymentMethod,
			order.CreatedAt,
			order.UpdatedAt,
		)

	if err != nil {
		return fmt.Errorf("обновлен заказ: %w", err)
	}

	return nil
}
