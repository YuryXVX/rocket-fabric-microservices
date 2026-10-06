package application

import (
	"context"
	"errors"
	errs "inventory/internal/errors"
	"inventory/internal/model/entity"
	"inventory/internal/service/input"
)

func (s *service) List(ctx context.Context, input input.PartFilter) ([]*entity.Part, error) {
	if len(input.UUIDs) > 0 {
		for _, f := range input.UUIDs {
			if !validateUUID(f.String()) {
				return []*entity.Part{}, errs.ErrInvalidUUID
			}
		}
	}

	list, err := s.inventoryRepository.List(ctx, input)

	if err != nil {
		if errors.Is(err, errs.ErrPartNotFound) {
			return []*entity.Part{}, errs.ErrPartNotFound
		}
		return nil, err
	}

	return list, nil
}
