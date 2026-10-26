package application

import (
	"context"
	"inventory/internal/model/entity"
	"inventory/internal/service/input"
)

type InventoryRepository interface {
	List(ctx context.Context, input input.PartFilter) ([]*entity.Part, error)
	ListForUpdate(ctx context.Context, input input.PartFilter) ([]*entity.Part, error)
	Get(ctx context.Context, uuid string) (*entity.Part, error)
	UpdateInventoryBatch(ctx context.Context, parts []*entity.Part) error
	CommitParts(ctx context.Context, input []string) error
}

type CompatibilityChecker interface {
	Check([]*entity.Part) error
}

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
