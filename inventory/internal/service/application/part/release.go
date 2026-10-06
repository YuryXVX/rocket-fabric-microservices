package application

import (
	"context"
	"fmt"
	"inventory/internal/service/input"
)

func (a *service) Release(ctx context.Context, input input.PartFilter) error {
	return a.txManager.Do(ctx, func(ctx context.Context) error {
		parts, err := a.List(ctx, input)

		if err != nil {
			return err
		}

		for _, p := range parts {
			p.ReleasePart()
		}

		fmt.Printf("part update release %+v", parts)

		return a.inventoryRepository.UpdateInventoryBatch(ctx, parts)
	})
}
