package v1

import (
	"context"
	"errors"
	"net/http"
	errs "order/internal/errors"
	"order/internal/model"

	orderv1 "shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	orderUUID, err := a.service.Pay(ctx, params.OrderUUID, model.PaymentMethod(req.PaymentMethod))

	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return &orderv1.PayOrderNotFound{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		}

		if errors.Is(err, errs.ErrOrderAlreadyPaid) {
			return &orderv1.PayOrderConflict{
				Code:    http.StatusConflict,
				Message: err.Error(),
			}, nil
		}

		return nil, err
	}

	return &orderv1.PayOrderResponse{
		TransactionUUID: orderUUID,
	}, nil
}
