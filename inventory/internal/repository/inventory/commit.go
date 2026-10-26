package inventory

import (
	"context"
	"fmt"
	errs "inventory/internal/errors"
)

func (r *repository) CommitParts(ctx context.Context, input []string) error {
	if len(input) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)

	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	defer tx.Rollback(ctx)

	selectQuery := `
		SELECT uuid, stock_quantity, reserved 
		FROM parts 
		WHERE uuid = ANY($1) 
		FOR UPDATE`

	uuids := make([]string, 0)

	for _, uuid := range input {
		if uuid == "" {
			continue
		}

		uuids = append(uuids, uuid)
	}

	rows, err := tx.Query(ctx, selectQuery, uuids)

	if err != nil {
		return fmt.Errorf("failed to select parts for update: %w", err)
	}

	defer rows.Close()

	foundCount := 0
	for rows.Next() {
		foundCount++
		var uuid string
		var stock, reserved int
		if err := rows.Scan(&uuid, &stock, &reserved); err != nil {
			return fmt.Errorf("failed to scan part: %w", err)
		}

		// Если хотя бы у одной детали не хватает остатков — сразу откат (FAILED_PRECONDITION)
		if stock <= 0 || reserved <= 0 {
			return errs.ErrPartsNotAvailable
		}
	}

	if rows.Err() != nil {
		return fmt.Errorf("rows error: %w", rows.Err())
	}

	// Если нашли меньше строк, чем передали UUID — значит каких-то деталей нет в базе (NOT_FOUND)
	if foundCount != len(uuids) {
		return errs.ErrPartNotFound
	}

	updateQuery := `
		UPDATE parts
		SET stock_quantity = stock_quantity - 1,
		    reserved = reserved - 1
		WHERE uuid = ANY($1)`

	_, err = tx.Exec(ctx, updateQuery, uuids)
	if err != nil {
		return fmt.Errorf("failed to update parts: %w", err)
	}

	return tx.Commit(ctx)
}
