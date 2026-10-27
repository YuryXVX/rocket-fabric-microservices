package order

import (
	"context"
	"fmt"
	"order/internal/model"

	"github.com/google/uuid"
)

func (s *service) CommitOrder(ctx context.Context, orderUUID uuid.UUID) error {
	return s.txManager.Do(ctx, func(ctx context.Context) error {
		order, err := s.orderRepository.Get(ctx, orderUUID)

		if err != nil {
			return err
		}

		if order.Status == model.OrderStatusAssembled {
			return fmt.Errorf("Заказ с uuid - %w уже в статусе ASSEMBLED", order.UUID.String())
		}

		uuids := make([]uuid.UUID, len(order.Items))

		for i, uuid := range order.Items {
			uuids[i] = uuid.PartUUID
		}

		err = s.inventoryClient.CommitParts(ctx, uuids)

		if err != nil {
			return err
		}

		order.Status = model.OrderStatusAssembled

		if err = s.orderRepository.Update(ctx, order); err != nil {
			return fmt.Errorf("Ошибка при обновлении статуса заказа")
		}

		return nil
	})
}
