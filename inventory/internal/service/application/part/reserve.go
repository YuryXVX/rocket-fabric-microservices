package application

import (
	"context"
	"inventory/internal/service/input"
)

func (a *service) Reserve(ctx context.Context, input input.PartFilter) error {
	return a.txManager.Do(ctx, func(ctx context.Context) error {

		parts, err := a.List(ctx, input)

		if err != nil {
			return err
		}

		for _, p := range parts {
			p.ReservedPart()
		}

		return a.inventoryRepository.UpdateInventoryBatch(ctx, parts)
	})
}
