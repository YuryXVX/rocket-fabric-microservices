package v1

import (
	"context"
	"errors"
	"net/http"
	"order/internal/api/converter"
	errs "order/internal/errors"
	orderv1 "shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	order, err := a.service.Create(ctx, converter.CreteRequestToInput(req))

	if err != nil {
		if errors.Is(err, errs.ErrPartNotFound) {
			return &orderv1.CreateOrderNotFound{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		}

		return nil, err
	}

	return &orderv1.CreateOrderResponse{
		OrderUUID:  order.UUID,
		TotalPrice: order.TotalPrice(),
	}, nil
}
