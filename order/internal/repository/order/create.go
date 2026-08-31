package order

import (
	"context"
	"order/internal/model"
)

func (r *repository) Create(ctx context.Context, order model.Order) error {
	return r.tx.Do(ctx, func(ctx context.Context) error {
		err := r.createOrder(ctx, order)

		if err != nil {
			return err
		}

		return r.createItems(ctx, order)
	})
}
