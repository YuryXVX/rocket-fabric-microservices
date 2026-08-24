package v1

import (
	"context"
	"order/internal/api/converter"
	orderv1 "shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	in := converter.CreteRequestToInput(req)

	a.service.Create(ctx, in)
	return &orderv1.CreateOrderResponse{}, nil
}
