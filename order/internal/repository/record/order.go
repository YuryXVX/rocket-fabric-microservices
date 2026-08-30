package record

import (
	"time"

	"github.com/google/uuid"
)

type OrderRecord struct {
	OrderUUID       uuid.UUID      `db:"uuid"`
	Status          OrderStatus    `db:"status"`
	TransactionUUID *uuid.UUID     `db:"transaction_uuid"`
	PaymentMethod   *PaymentMethod `db:"payment_method"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       *time.Time     `db:"updated_at"` // Добавлено из SQL-схемы
}

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type PaymentMethod string

const (
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)
