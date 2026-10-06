package inventory

import (
	"context"
	"database/sql"
	"errors"

	errs "inventory/internal/errors"
	"inventory/internal/model/entity"
	"inventory/internal/repository/converter"
	"inventory/internal/repository/record"
)

const getPartQuery = `
	SELECT 
		uuid, name, description, part_type, 
		price, stock_quantity, created_at, updated_at 
	FROM parts 
	WHERE uuid = $1
`

func (r *repository) Get(ctx context.Context, uuid string) (*entity.Part, error) {
	var record record.Part

	row := r.pool.QueryRow(ctx, getPartQuery, uuid)

	err := row.Scan(
		&record.UUID,
		&record.Name,
		&record.Description,
		&record.PartType,
		&record.Price,
		&record.StockQuantity,
		&record.CreatedAt,
		&record.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPartNotFound
		}
		return nil, err
	}

	part, err := converter.RecordPartToDomain(record)

	if err != nil {
		return nil, err
	}

	return &part, nil
}
