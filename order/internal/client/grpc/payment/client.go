package payment

import (
	"context"
	v1 "shared/pkg/proto/payment/v1"
)

type paymentClient struct {
	grpc v1.BillingServiceClient
}

func New(c v1.BillingServiceClient) *paymentClient {
	return &paymentClient{
		grpc: c,
	}
}

func (p *paymentClient) PayOrder(ctx context.Context) {
}
