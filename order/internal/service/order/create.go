package order

import (
	"context"
	"fmt"
	"order/internal/model"
	"order/internal/service/input"
)

func (s *service) Create(ctx context.Context, in *input.CreateOrderInput) (model.Order, error) {
	parts, err := s.inventoryClient.ListParts(ctx, in.PartUUIDs())

	if err != nil {
		return model.Order{}, err
	}

	for _, p := range parts {
		fmt.Println(p)
	}

	return model.Order{}, nil
}
