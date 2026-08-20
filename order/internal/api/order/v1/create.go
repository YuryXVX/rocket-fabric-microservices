package v1

import (
	"context"

	orderv1 "shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	return &orderv1.CreateOrderResponse{}, nil
}
