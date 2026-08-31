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
	items, err := s.inventoryClient.ListParts(ctx, in.PartUUIDs())

	if err != nil {
		return model.Order{}, fmt.Errorf("получение детали %s", err.Error())
	}

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

	if err := s.orderRepository.Create(ctx, order); err != nil {
		return model.Order{}, fmt.Errorf("сохранение заказа %s", order.UUID)

	}

	return order, nil
}
