package application

import (
	"context"
	errs "inventory/internal/errors"
	"inventory/internal/model/entity"

	"inventory/internal/service/input"

	"github.com/google/uuid"
)

func (s *service) ValidateCompatibility(ctx context.Context, input input.PartFilter) error {
	resolved, err := s.resolvedShipSlots(ctx, input)

	if err != nil {
		return err
	}

	return s.compatibilityChecker.Check(resolved)
}

func (s *service) resolvedShipSlots(ctx context.Context, input input.PartFilter) ([]*entity.Part, error) {
	seen := make(map[uuid.UUID]struct{}, len(input.UUIDs))

	for _, uuid := range input.UUIDs {
		if _, ok := seen[uuid]; ok {
			return nil, errs.ErrDuplicatePartUUID
		}

		seen[uuid] = struct{}{}
	}

	return s.List(ctx, input)
}
