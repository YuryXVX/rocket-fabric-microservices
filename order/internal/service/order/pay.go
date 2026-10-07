package order

import (
	"context"
	"errors"
	"fmt"
	errs "order/internal/errors"
	"order/internal/model"
	"time"

	"github.com/google/uuid"
)

func (s *service) Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	var transactionUUID uuid.UUID

	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		order, err := s.orderRepository.Get(ctx, orderUUID)

		if err != nil {
			if errors.Is(err, errs.ErrOrderNotFound) {
				return fmt.Errorf("заказ с %s не найден", orderUUID)
			}

			return err
		}

		if order.Status == model.OrderStatusPaid {
			return errs.ErrOrderAlreadyPaid
		}

		if order.Status != model.OrderStatusPendingPayment {
			return fmt.Errorf("заказ с uuid %s находится не в статусе ожидания оплаты", orderUUID)
		}

		order.Status = model.OrderStatusPaid
		order.PaymentMethod = &method
		order.TransactionUUID = &transactionUUID
		order.UpdatedAt = time.Now()

		if err = s.orderRepository.Update(ctx, order); err != nil {
			return fmt.Errorf("при оплате %s", orderUUID)
		}

		uuid, err := s.paymentClient.PayOrder(ctx, orderUUID.String(), model.PaymentMethodCard)

		if err != nil {
			return err
		}

		transactionUUID = uuid

		return nil
	})

	if err != nil {
		return uuid.UUID{}, err
	}

	return transactionUUID, nil
}
