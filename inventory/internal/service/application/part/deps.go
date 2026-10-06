package application

import (
	"context"
	"inventory/internal/model/entity"
	"inventory/internal/service/input"
)

type InventoryRepository interface {
	List(ctx context.Context, input input.PartFilter) ([]*entity.Part, error)
	Get(ctx context.Context, uuid string) (*entity.Part, error)
}

type CompatibilityChecker interface {
	Check([]*entity.Part) error
}
