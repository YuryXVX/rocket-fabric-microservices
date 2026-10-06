package inventory

import (
	"context"
	"fmt"
	"inventory/internal/model/entity"

	"github.com/jackc/pgx/v5"
)

const query = "UPDATE parts SET reserved = $1 WHERE uuid = $2"

func (r *repository) UpdateInventoryBatch(ctx context.Context, parts []*entity.Part) error {

	batch := &pgx.Batch{}

	for _, part := range parts {
		batch.Queue(
			query,
			part.Reserved(),
			part.UUID().String(),
		)
	}

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < len(parts); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("ошибка обновления в пакете на шаге %d: %w", i, err)
		}
	}

	return nil
}
