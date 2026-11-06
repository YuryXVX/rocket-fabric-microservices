package v1

import (
	"context"
	"errors"
	"net/http"
	"order/internal/api/converter"
	errs "order/internal/errors"
	"plaform/pkg/auth"
	orderv1 "shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	userUUID, ok := auth.UserUUIDFromContext(ctx)

	if !ok {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "uuid пользователя отсутсвует или невалидный",
		}, nil
	}

	order, err := a.service.Create(ctx, converter.CreteRequestToInput(req, userUUID))

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
