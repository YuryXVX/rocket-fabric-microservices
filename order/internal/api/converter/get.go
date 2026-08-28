package converter

import (
	"order/internal/model"
	orderv1 "shared/pkg/openapi/order/v1"

	"github.com/google/uuid"
)

func NewOptUUID(val *uuid.UUID) orderv1.OptNilUUID {
	if val == nil {
		return orderv1.OptNilUUID{
			Set:  true,
			Null: true,
		}
	}

	return orderv1.OptNilUUID{
		Value: *val,
		Set:   true,
		Null:  false,
	}
}

func covertModelPaymentMethodToResponse(method *model.PaymentMethod) orderv1.OptNilPaymentMethod {
	if method == nil {
		return orderv1.NewOptNilPaymentMethod("")
	}

	var m orderv1.PaymentMethod

	switch *method {
	case model.PaymentMethodCard:
		m = orderv1.PaymentMethodCARD
	case model.PaymentMethodSBP:
		m = orderv1.PaymentMethodSBP
	case model.PaymentMethodCreditCard:
		m = orderv1.PaymentMethodCREDITCARD
	case model.PaymentMethodInvestorMoney:
		m = orderv1.PaymentMethodINVESTORMONEY
	}

	return orderv1.NewOptNilPaymentMethod(m)
}

func OrderModelToResponse(order model.Order) *orderv1.OrderDto {
	response := orderv1.OrderDto{
		OrderUUID:     order.UUID,
		Status:        orderv1.OrderStatus(order.Status),
		HullUUID:      order.HullUUID(),
		EngineUUID:    order.EngineUUID(),
		TotalPrice:    order.TotalPrice(),
		CreatedAt:     order.CreatedAt,
		PaymentMethod: covertModelPaymentMethodToResponse(order.PaymentMethod),
	}

	if order.WeaponUUID() != nil {
		response.WeaponUUID = NewOptUUID(order.WeaponUUID())
	}

	if order.ShieldUUID() != nil {
		response.ShieldUUID = NewOptUUID(order.ShieldUUID())
	}

	if order.TransactionUUID != nil {
		response.TransactionUUID = NewOptUUID(order.TransactionUUID)
	}

	if order.PaymentMethod != nil {
		response.PaymentMethod = covertModelPaymentMethodToResponse(order.PaymentMethod)
	}

	return &response
}
