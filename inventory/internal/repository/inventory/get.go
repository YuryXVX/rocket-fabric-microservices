package inventory

import (
	"context"

	errs "inventory/internal/errors"
	"inventory/internal/model"
	"inventory/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, uuid string) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.parts[uuid]

	if !ok {
		return nil, errs.ErrPartNotFound
	}

	return converter.RecordPartToModel(part), nil
}
