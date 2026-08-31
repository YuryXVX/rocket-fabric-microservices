package order

import (
	"context"
	"fmt"
	"order/internal/model"
	"order/internal/repository/converter"
	"order/internal/repository/record"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

var fields = []string{
	"o.uuid",
	"o.status",
	"o.transaction_uuid",
	"o.payment_method",
	"o.created_at",
	"o.updated_at",
	"oi.part_uuid",
	"oi.part_type",
	"oi.price",
}

func (r *repository) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	var rec record.OrderRecord

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := builder.Select(fields...).
		From("orders o").
		LeftJoin("order_items oi ON o.uuid = oi.order_uuid").
		Where(squirrel.Eq{"o.uuid": orderUUID}).
		ToSql()

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).
		Query(ctx, sql, args...)

	if err != nil {
		return model.Order{}, fmt.Errorf("failed to build sql: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var partUUID *uuid.UUID
		var partType *string
		var price *int64

		// Сканируем ВСЕ поля из Select
		err = rows.Scan(
			&rec.OrderUUID, &rec.Status, &rec.TransactionUUID, &rec.PaymentMethod, &rec.CreatedAt, &rec.UpdatedAt,
			&partUUID, &partType, &price,
		)
		if err != nil {
			return model.Order{}, fmt.Errorf("failed to scan row: %w", err)
		}

		if partUUID != nil && partType != nil && price != nil {
			fmt.Println(*partType, *partUUID, *price)
			rec.Items = append(rec.Items, record.PartRecord{
				OrderUUID: rec.OrderUUID,
				PartUUID:  *partUUID,
				PartType:  record.PartType(*partType),
				Price:     *price,
			})
		}
	}

	return converter.RecordOrderToModel(rec), nil
}
