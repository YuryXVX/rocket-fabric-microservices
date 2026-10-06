package application

import (
	"context"
	errs "inventory/internal/errors"
	"inventory/internal/model/entity"
)

func (s *service) Get(ctx context.Context, uuid string) (*entity.Part, error) {
	if !validateUUID(uuid) {
		return &entity.Part{}, errs.ErrInvalidUUID
	}

	return s.inventoryRepository.Get(ctx, uuid)
}
