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
	order, err := s.orderRepository.Get(ctx, orderUUID)

	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return uuid.UUID{}, fmt.Errorf("заказ с %s не найден", orderUUID)
		}

		return uuid.UUID{}, err
	}

	if order.Status == model.OrderStatusPaid {
		return uuid.UUID{}, errs.ErrOrderAlreadyPaid
	}

	if order.Status != model.OrderStatusPendingPayment {
		return uuid.UUID{}, fmt.Errorf("заказ с uuid %s находится не в статусе ожидания оплаты", orderUUID)
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, orderUUID.String(), model.PaymentMethodCard)

	if err != nil {
		return uuid.UUID{}, nil
	}

	order.Status = model.OrderStatusPaid
	order.PaymentMethod = &method
	order.TransactionUUID = &transactionUUID
	order.UpdatedAt = time.Now()

	if err = s.orderRepository.Update(ctx, order); err != nil {
		return uuid.UUID{}, fmt.Errorf("при оплате %s", orderUUID)
	}

	return transactionUUID, nil
}
