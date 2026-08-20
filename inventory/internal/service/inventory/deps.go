package inventory

import (
	"context"

	"inventory/internal/model"
	"inventory/internal/service/input"
)

type InventoryRepository interface {
	List(ctx context.Context, input input.PartFilter) ([]*model.Part, error)
	Get(ctx context.Context, uuid string) (*model.Part, error)
}
