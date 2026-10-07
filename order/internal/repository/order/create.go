package order

import (
	"context"
	"order/internal/model"
)

func (r *repository) Create(ctx context.Context, order model.Order) error {
	return r.createOrder(ctx, order)
}
