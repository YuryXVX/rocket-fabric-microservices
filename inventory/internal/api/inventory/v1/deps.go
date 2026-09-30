package v1

import (
	"context"

	"inventory/internal/model"
	"inventory/internal/model/entity"
	"inventory/internal/service/input"
)

type ServiceInventory interface {
	Get(ctx context.Context, uuid string) (*model.Part, error)
	List(ctx context.Context, input input.PartFilter) ([]*entity.Part, error)
}
