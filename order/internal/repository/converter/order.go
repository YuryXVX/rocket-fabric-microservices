package converter

import (
	"order/internal/model"
	"order/internal/repository/record"
)

func OrderModelToRecord(in *model.Order) record.OrderRecord {
	return record.OrderRecord{
		OrderUUID:       in.UUID,
		CreatedAt:       in.CreatedAt,
		TransactionUUID: in.TransactionUUID,
		Status:          orderStatusToRecordStatus(in.Status),
		PaymentMethod:   paymentMethodToRecord(in.PaymentMethod),
	}
}

func RecordOrderToModel(in record.OrderRecord) model.Order {
	m := model.Order{
		UUID:            in.OrderUUID,
		Status:          recordStatusToOrderStatus(in.Status),
		CreatedAt:       in.CreatedAt,
		TransactionUUID: in.TransactionUUID,
		PaymentMethod:   recordPaymentMethodToModel(in.PaymentMethod),
	}

	if len(in.Items) > 0 {
		for _, item := range in.Items {
			m.Items = append(m.Items, RecordItemToModel(item))

		}
	}

	return m
}

func orderStatusToRecordStatus(in model.OrderStatus) record.OrderStatus {
	var s record.OrderStatus

	switch in {
	case model.OrderStatusPendingPayment:
		s = record.OrderStatusPendingPayment
	case model.OrderStatusCancelled:
		s = record.OrderStatusCancelled
	case model.OrderStatusPaid:
		s = record.OrderStatusPaid
	}

	return s
}

func recordStatusToOrderStatus(in record.OrderStatus) model.OrderStatus {
	var s model.OrderStatus

	switch in {
	case record.OrderStatusPendingPayment:
		s = model.OrderStatusPendingPayment
	case record.OrderStatusCancelled:
		s = model.OrderStatusCancelled
	case record.OrderStatusPaid:
		s = model.OrderStatusPaid
	}

	return s
}

func recordPaymentMethodToModel(in *record.PaymentMethod) *model.PaymentMethod {
	if in == nil {
		return nil
	}

	var out model.PaymentMethod
	switch *in {
	case record.PaymentMethodCard:
		out = model.PaymentMethodCard
	case record.PaymentMethodSBP:
		out = model.PaymentMethodSBP
	case record.PaymentMethodCreditCard:
		out = model.PaymentMethodCreditCard
	case record.PaymentMethodInvestorMoney:
		out = model.PaymentMethodInvestorMoney
	}

	return &out
}

func paymentMethodToRecord(in *model.PaymentMethod) *record.PaymentMethod {
	if in == nil {
		return nil
	}

	var out record.PaymentMethod
	switch *in {
	case model.PaymentMethodCard:
		out = record.PaymentMethodCard
	case model.PaymentMethodSBP:
		out = record.PaymentMethodSBP
	case model.PaymentMethodCreditCard:
		out = record.PaymentMethodCreditCard
	case model.PaymentMethodInvestorMoney:
		out = record.PaymentMethodInvestorMoney
	}

	return &out
}
