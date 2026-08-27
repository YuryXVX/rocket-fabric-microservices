package part

import (
	"context"
	"order/internal/model"
	"order/internal/repository/converter"

	"github.com/google/uuid"
)

func (r *repository) SavePart(ctx context.Context, orderUUID uuid.UUID, partRecord []model.OrderItem) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.store[orderUUID] = converter.OrderItemsToRecords(partRecord)

	return nil
}
