package converter

import (
	"github.com/google/uuid"
	errs "payment/internal/errors"
	"payment/internal/model"
	"payment/internal/service/input"
	v1 "shared/pkg/proto/payment/v1"
)

func ValidUuid(req *v1.PayOrderRequest) error {
	if req.OrderUuid == "" {
		return errs.ErrInvalidOrderUUID
	}

	if _, err := uuid.Parse(req.OrderUuid); err != nil {
		return errs.ErrInvalidOrderUUID
	}

	return nil
}

func RequestPayToInput(req *v1.PayOrderRequest) *input.PayOrderInput {
	if req.OrderUuid == "" {
		return nil
	}

	return &input.PayOrderInput{
		OrderUUID:     uuid.MustParse(req.OrderUuid),
		PaymentMethod: convertPaymentMethod(req.PaymentMethod),
	}
}

func convertPaymentMethod(in v1.PaymentMethod) model.PaymentMethod {
	switch in {
	case v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	case v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case v1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case v1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	default:
		return model.PaymentMethodUnspecified
	}
}
