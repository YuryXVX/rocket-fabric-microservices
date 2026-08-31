package order

import (
	"context"
	"order/internal/model"
)

func (r *repository) Update(ctx context.Context, order model.Order) error {
	return r.tx.Do(ctx, func(ctx context.Context) error {
		return r.updateOrder(ctx, order)
	})
}
