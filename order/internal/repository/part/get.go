package part

import (
	"context"
	"fmt"
	"order/internal/model"

	"github.com/google/uuid"
)

func (r *repository) GetPart(ctx context.Context, orderUUID uuid.UUID) ([]model.OrderItem, error) {
	var orderItems []model.OrderItem

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(
		ctx,
		"SELECT * FROM order_items WHERE order_uuid=($1)",
		orderUUID,
	)

	if err != nil {
		return orderItems, fmt.Errorf("query order failed: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var o model.OrderItem

		err := rows.Scan(&orderUUID, &o.PartUUID, &o.PartType, &o.Price)

		if err != nil {
			return nil, fmt.Errorf("scan order failed: %w", err)
		}

		orderItems = append(orderItems, o)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("итерация завершилась с ошибкой: %w", err)
	}

	return orderItems, nil
}
