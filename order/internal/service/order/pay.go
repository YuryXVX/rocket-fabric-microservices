package order

import (
	"context"
	"errors"
	"fmt"
	errs "order/internal/errors"
	"order/internal/model"

	"github.com/google/uuid"
)

func (s *service) Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	_, err := s.orderRepository.Get(ctx, orderUUID)

	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return uuid.UUID{}, fmt.Errorf("заказ с %s не найден", orderUUID)
		}

		return uuid.UUID{}, err
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, orderUUID.String(), model.PaymentMethodCard)

	if err != nil {
		return uuid.UUID{}, err
	}

	return transactionUUID, nil
}
