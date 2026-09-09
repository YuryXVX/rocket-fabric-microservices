package app

import (
	"context"
	apiV1 "payment/internal/api/payment/v1"
	service "payment/internal/service/payment"
	v1 "shared/pkg/proto/payment/v1"
)

type diContainer struct {
	// service
	paymentService apiV1.PayService
	// handlers
	paymentHandlers v1.BillingServiceServer
}

func (di *diContainer) PaymentService(_ context.Context) apiV1.PayService {
	if di.paymentService == nil {
		di.paymentService = service.New()
	}

	return di.paymentService
}

func (di *diContainer) PaymentHandlers(ctx context.Context) v1.BillingServiceServer {
	if di.paymentHandlers == nil {
		di.paymentHandlers = apiV1.New(di.PaymentService(ctx))
	}

	return di.paymentHandlers
}
