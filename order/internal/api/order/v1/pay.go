package v1

import (
	"context"
	"order/internal/model"

	orderv1 "shared/pkg/openapi/order/v1"
)

// Проводит оплату заказа.
func (a *api) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	orderUUID, err := a.service.Pay(ctx, params.OrderUUID, model.PaymentMethod(req.PaymentMethod))

	if err != nil {
		return nil, err
	}

	return &orderv1.PayOrderResponse{
		TransactionUUID: orderUUID,
	}, nil
}
