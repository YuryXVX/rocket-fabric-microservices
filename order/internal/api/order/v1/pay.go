package v1

import (
	"context"

	"github.com/google/uuid"
	orderv1 "shared/pkg/openapi/order/v1"
)

// Проводит оплату заказа.
func (a *api) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	uuid := uuid.New()

	return &orderv1.PayOrderResponse{
		TransactionUUID: uuid,
	}, nil
}
