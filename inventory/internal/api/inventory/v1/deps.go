package v1

import (
	"context"

	"inventory/internal/model/entity"
	"inventory/internal/service/input"
)

type ApplicationService interface {
	Get(context.Context, string) (*entity.Part, error)
	List(context.Context, input.PartFilter) ([]*entity.Part, error)
	Release(context.Context, input.PartFilter) error
	Reserve(context.Context, input.PartFilter) error
	ValidateCompatibility(context.Context, input.PartFilter) error
	CommitParts(context.Context, *input.PartFilter) error
}
