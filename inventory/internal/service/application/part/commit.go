package application

import (
	"context"
	errs "inventory/internal/errors"
	"inventory/internal/service/input"
)

func (s *service) CommitParts(ctx context.Context, input *input.PartFilter) error {
	if len(input.UUIDs) > 0 {
		for _, f := range input.UUIDs {
			if !validateUUID(f.String()) {
				return errs.ErrInvalidUUID
			}
		}
	}

	uuids := make([]string, len(input.UUIDs))

	for _, uuid := range input.UUIDs {
		uuids = append(uuids, uuid.String())
	}

	return s.inventoryRepository.CommitParts(ctx, uuids)
}
