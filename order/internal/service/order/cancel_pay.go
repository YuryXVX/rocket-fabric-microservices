package order

import (
	"context"
	"fmt"
	"order/internal/model"

	"github.com/google/uuid"
)

func (s *service) Cancel(ctx context.Context, orderUUID uuid.UUID) error {
	order, err := s.orderRepository.Get(ctx, orderUUID)

	if err != nil {
		return fmt.Errorf("заказ с uuid %s отсутствует", orderUUID)
	}

	if order.Status != model.OrderStatusPendingPayment {
		return fmt.Errorf("заказ с uuid %s находится не в статусе ожидания оплаты", orderUUID)
	}

	order.Status = model.OrderStatusCancelled

	if err := s.orderRepository.Update(ctx, order); err != nil {
		return fmt.Errorf("Ошибка при отмене оплаты")
	}

	return nil
}
