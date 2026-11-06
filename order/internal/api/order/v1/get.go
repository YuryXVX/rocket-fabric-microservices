package v1

import (
	"context"
	"errors"
	"net/http"
	"order/internal/api/converter"
	errs "order/internal/errors"

	orderV1 "shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	order, err := a.service.Get(ctx, params.OrderUUID)

	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return &orderV1.GetOrderNotFound{
				Code: http.StatusNotFound,
			}, nil
		}

		return &orderV1.GetOrderInternalServerError{}, nil
	}

	return converter.OrderModelToResponse(order), nil
}
