package order

import (
	"context"
	"fmt"
	"order/internal/model"

	"github.com/google/uuid"
)

func (s *service) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)

	if err != nil {
		fmt.Println(err.Error())
		return model.Order{}, err
	}

	return order, nil
}
