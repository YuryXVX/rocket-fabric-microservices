package order

import (
	"context"
	"order/internal/model"
)

func (r *repository) Update(ctx context.Context, order model.Order) error {
	return r.updateOrder(ctx, order)
}
