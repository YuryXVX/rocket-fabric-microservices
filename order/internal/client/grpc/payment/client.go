package payment

import (
	"context"
	"order/internal/client/grpc/payment/converter"
	"order/internal/model"
	v1 "shared/pkg/proto/payment/v1"

	"github.com/google/uuid"
)

type paymentClient struct {
	grpc v1.BillingServiceClient
}

func New(c v1.BillingServiceClient) *paymentClient {
	return &paymentClient{
		grpc: c,
	}
}

func (p *paymentClient) PayOrder(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
	response, err := p.grpc.PayOrder(ctx, converter.InputToPaymentRequest(orderUUID, paymentMethod))

	if err != nil {
		return uuid.UUID{}, err
	}

	return uuid.MustParse(response.TransactionUuid), nil
}
