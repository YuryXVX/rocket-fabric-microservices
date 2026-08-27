package converter

import (
	"order/internal/model"
	v1 "shared/pkg/proto/payment/v1"
)

func InputToPaymentRequest(orderUUID string, paymentMethod model.PaymentMethod) *v1.PayOrderRequest {
	return &v1.PayOrderRequest{
		OrderUuid:     orderUUID,
		PaymentMethod: toRequestPaymentMethod(paymentMethod),
	}
}

func toRequestPaymentMethod(in model.PaymentMethod) v1.PaymentMethod {
	var m v1.PaymentMethod

	switch in {
	case model.PaymentMethodCard:
		m = v1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodCreditCard:
		m = v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodSBP:
		m = v1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodInvestorMoney:
		m = v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		m = v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}

	return m
}
