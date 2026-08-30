package inventory

import (
	"context"
	"fmt"

	errs "inventory/internal/errors"
	"inventory/internal/model"
	"inventory/internal/repository/converter"
	"inventory/internal/repository/record"
	"inventory/internal/service/input"
)

func (r *repository) List(ctx context.Context, input input.PartFilter) ([]*model.Part, error) {
	query := "SELECT uuid, name, description, part_type, price, stock_quantity, created_at, updated_at FROM parts"

	var args []interface{}

	if len(input.UUIDs) > 0 {
		query += " WHERE uuid = ANY($1)"
		args = append(args, input.UUIDs)
	} else if input.PartType != "" {
		query += " WHERE part_type = $1"
		args = append(args, input.PartType)
	}

	fmt.Println(args...)

	query += " ORDER BY LOWER(name) ASC"

	rows, err := r.pool.Query(ctx, query, args...)

	if err != nil {
		return nil, fmt.Errorf("failed to query parts: %w", err)
	}

	defer rows.Close()

	var parts []*model.Part

	for rows.Next() {
		var p record.Part

		var partTypeStr string

		err := rows.Scan(
			&p.UUID,
			&p.Name,
			&p.Description,
			&partTypeStr,
			&p.Price,
			&p.StockQuantity,
			&p.CreatedAt,
			&p.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan part: %w", err)
		}

		p.PartType = converter.ConvertStringToPartType(partTypeStr)

		parts = append(parts, converter.RecordPartToModel(p))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	if len(input.UUIDs) > 0 && len(input.UUIDs) != len(parts) {
		return []*model.Part{}, errs.ErrPartNotFound
	}

	return parts, nil
}
