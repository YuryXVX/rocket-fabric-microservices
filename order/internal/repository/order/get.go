package order

import (
	"context"
	errs "order/internal/errors"
	"order/internal/model"
	"order/internal/repository/converter"

	"github.com/google/uuid"
)

func (r *repository) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.store[orderUUID]

	if !ok {
		return model.Order{}, errs.ErrOrderNotFound
	}

	return converter.RecordOrderToModel(order), nil
}
