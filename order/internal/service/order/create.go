package order

import (
	"context"
	"fmt"
	errs "order/internal/errors"
	"order/internal/model"
	"order/internal/service/input"
	"plaform/pkg/auth"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func (s *service) withGrpcAuth(ctx context.Context) context.Context {
	token, ok := auth.SessionUUIDFromContext(ctx)
	if !ok {
		return ctx
	}

	return metadata.AppendToOutgoingContext(ctx, auth.SessionTokenKey, token)
}

func (s *service) Create(ctx context.Context, in *input.CreateOrderInput) (model.Order, error) {
	var order model.Order

	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		ctx = s.withGrpcAuth(ctx)

		items, err := s.inventoryClient.ListParts(ctx, in.PartUUIDs())

		if err != nil {
			return fmt.Errorf("получение детали %s", err.Error())
		}

		if len(items) == 0 || len(items) != len(in.PartUUIDs()) {
			return errs.ErrPartNotFound
		}

		err = s.inventoryClient.ValidateCompatibility(ctx, *in)

		if err != nil {
			return fmt.Errorf("Заказ не прошле проверку [ValidateCompatibility] %s ", err.Error())
		}

		err = s.inventoryClient.ReserveParts(ctx, in.PartUUIDs())

		if err != nil {
			return fmt.Errorf("Ошибка при резервировании заказ [ReserveParts] %s ", err.Error())
		}

		order = model.Order{
			UUID:            uuid.New(),
			UserUUID:        in.UserUUID,
			Items:           items,
			TransactionUUID: nil,
			PaymentMethod:   nil,
			Status:          model.OrderStatusPendingPayment,
			CreatedAt:       time.Now(),
		}

		err = s.orderRepository.Create(ctx, order)

		if err != nil {
			return fmt.Errorf("oшибка при записи заказа")
		}

		err = s.orderRepository.CreateItems(ctx, order)

		if err != nil {
			return fmt.Errorf("oшибка при записи деталей")
		}

		return nil
	})

	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
