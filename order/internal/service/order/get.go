package order

import (
	"context"
	"order/internal/model"

	"github.com/google/uuid"
)

func (s *service) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	items, err := s.partRepository.GetPart(ctx, orderUUID)

	if err != nil {
		return model.Order{}, err
	}

	order, err := s.orderRepository.Get(ctx, orderUUID)

	if err != nil {
		return model.Order{}, err
	}

	order.OrderItems(items)

	return order, nil
}
