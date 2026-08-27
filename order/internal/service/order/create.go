package order

import (
	"context"
	"fmt"
	errs "order/internal/errors"
	"order/internal/model"
	"order/internal/service/input"
	"time"

	"github.com/google/uuid"
)

func (s *service) Create(ctx context.Context, in *input.CreateOrderInput) (model.Order, error) {
	items, _ := s.inventoryClient.ListParts(ctx, in.PartUUIDs())

	if len(items) == 0 || len(items) != len(in.PartUUIDs()) {
		return model.Order{}, errs.ErrPartNotFound
	}

	order := model.Order{
		UUID:            uuid.New(),
		Items:           items,
		TransactionUUID: nil,
		PaymentMethod:   nil,
		Status:          model.OrderStatusPendingPayment,
		CreatedAt:       time.Now(),
	}

	err := s.orderRepository.Create(ctx, order)

	if err != nil {
		return model.Order{}, err
	}

	err = s.partRepository.SavePart(ctx, order.UUID, order.Items)

	if err != nil {
		return model.Order{}, fmt.Errorf("произошла ошибка при сохранении part %s", order.UUID)
	}

	return order, nil
}
