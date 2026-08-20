package v1

import (
	"context"

	"github.com/google/uuid"
	orderv1 "shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	return &orderv1.OrderDto{
		OrderUUID: uuid.New(),
		HullUUID:  uuid.New(),
	}, nil
}
