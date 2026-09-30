package inventory

import (
	"context"
	"errors"

	errs "inventory/internal/errors"
	"inventory/internal/model"
	"inventory/internal/model/entity"
	"inventory/internal/service/input"

	"github.com/google/uuid"
)

type service struct {
	inventoryRepository  InventoryRepository
	compatibilityChecker CompatibilityChecker
}

func NewService(repo InventoryRepository, checker CompatibilityChecker) *service {
	return &service{
		inventoryRepository:  repo,
		compatibilityChecker: checker,
	}
}

func (s *service) Get(ctx context.Context, uuid string) (*model.Part, error) {
	if !validateUUID(uuid) {
		return &model.Part{}, errs.ErrInvalidUUID
	}

	return s.inventoryRepository.Get(ctx, uuid)
}

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

func validateUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
