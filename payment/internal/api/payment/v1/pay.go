package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"payment/internal/api/converter"
	errs "payment/internal/errors"
	v1 "shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *v1.PayOrderRequest) (*v1.PayOrderResponse, error) {
	err := converter.ValidUuid(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка в uuid %s", err)
	}

	input := converter.RequestPayToInput(req)

	if input == nil {
		return nil, status.Errorf(codes.InvalidArgument, "неверный формат запроса")
	}

	uuid, err := a.service.Pay(ctx, *input)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidPaymentMethod) {
			return nil, status.Errorf(codes.InvalidArgument, "неверный метод оплаты")
		}

		return nil, err
	}

	return &v1.PayOrderResponse{
		TransactionUuid: uuid.String(),
	}, nil
}
