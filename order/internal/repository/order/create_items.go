package order

import (
	"context"
	"fmt"
	"order/internal/model"
	"strings"
)

func (r *repository) createItems(ctx context.Context, order model.Order) error {
	if len(order.Items) == 0 {
		return nil
	}

	const columnsCount = 4

	args := make([]interface{}, 0, len(order.Items)*columnsCount)

	var valueStrings strings.Builder

	fmt.Printf("order.Items: %v\n", order.Items)

	valueStrings.Grow(len(order.Items) * 20)

	for i, item := range order.Items {
		if i > 0 {
			valueStrings.WriteString(", ")
		}

		p := i * columnsCount
		fmt.Fprintf(&valueStrings, "($%d, $%d, $%d, $%d)", p+1, p+2, p+3, p+4)

		args = append(args, order.UUID, item.PartUUID, item.PartType, item.Price)
	}

	query := fmt.Sprintf(
		"INSERT INTO order_items (order_uuid, part_uuid, part_type, price) VALUES %s",
		valueStrings.String(),
	)

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("bulk insert order items failed: %w", err)
	}

	return nil
}
