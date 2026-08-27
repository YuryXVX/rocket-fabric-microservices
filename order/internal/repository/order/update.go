package order

import (
	"context"
	"order/internal/model"
	"order/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store[order.UUID] = converter.OrderModelToRecord(&order)

	return nil
}
