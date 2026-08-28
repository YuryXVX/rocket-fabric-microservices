package v1

import (
	"context"
	"errors"
	"net/http"
	errs "order/internal/errors"

	orderv1 "shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	err := a.service.Cancel(ctx, params.OrderUUID)

	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return &orderv1.CancelOrderNotFound{
				Code: http.StatusNotFound,
			}, nil
		}

		return nil, err
	}

	return &orderv1.CancelOrderResponse{}, nil
}
