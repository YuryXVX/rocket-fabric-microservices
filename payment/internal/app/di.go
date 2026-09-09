package app

import (
	"context"
	apiV1 "payment/internal/api/payment/v1"
	service "payment/internal/service/payment"
	"plaform/pkg/di"
	v1 "shared/pkg/proto/payment/v1"
)

type diContainer struct {
	// service
	paymentService di.Value[apiV1.PayService]
	// handlers
	paymentHandlers di.Value[v1.BillingServiceServer]
}

func (di *diContainer) PaymentService(ctx context.Context) apiV1.PayService {
	return di.paymentService.Get(ctx, func(ctx context.Context) apiV1.PayService {
		return service.New()
	})
}

func (di *diContainer) PaymentHandlers(ctx context.Context) v1.BillingServiceServer {
	return di.paymentHandlers.Get(ctx, func(ctx context.Context) v1.BillingServiceServer {
		return apiV1.New(di.PaymentService(ctx))
	})
}
