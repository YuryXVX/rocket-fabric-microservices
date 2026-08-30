package inventory

import (
	"context"
	"fmt"

	"inventory/internal/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (*model.Part, error) {
	var part model.Part

	row := r.pool.QueryRow(ctx, "SELECT * FROM parts WHERE uuid = $1", uuid)

	err := row.Scan(
		&part.UUID,
		&part.Name,
		&part.Description,
		&part.PartType,
		&part.Price,
		&part.StockQuantity,
		&part.CreatedAt,
		&part.UpdatedAt,
	)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &part, nil
}
