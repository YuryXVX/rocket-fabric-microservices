package part

import (
	"context"
	"order/internal/model"
	"order/internal/repository/converter"

	"github.com/google/uuid"
)

func (r *repository) GetPart(ctx context.Context, orderUUID uuid.UUID) ([]model.OrderItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items, ok := r.store[orderUUID]

	if !ok {
		return []model.OrderItem{}, nil
	}

	return converter.RecordItemsToModel(items), nil
}
