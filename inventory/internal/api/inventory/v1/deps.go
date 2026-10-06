package v1

import (
	"context"

	"inventory/internal/model/entity"
	"inventory/internal/service/input"
)

type ApplicationService interface {
	Get(ctx context.Context, uuid string) (*entity.Part, error)
	List(ctx context.Context, input input.PartFilter) ([]*entity.Part, error)
	Release(ctx context.Context, input input.PartFilter) error
	Reserve(ctx context.Context, input input.PartFilter) error
	ValidateCompatibility(ctx context.Context, input input.PartFilter) error
}
