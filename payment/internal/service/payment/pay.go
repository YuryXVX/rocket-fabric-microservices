package payment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	errs "payment/internal/errors"
	"payment/internal/service/input"
)

func (s *service) Pay(ctx context.Context, in input.PayOrderInput) (uuid.UUID, error) {
	if !in.PaymentMethod.IsValid() {
		return uuid.Nil, errs.ErrInvalidPaymentMethod
	}

	// 2. Генерация transaction_uuid
	transactionUUID := uuid.New()

	// 3. Логирование
	slog.InfoContext(ctx, "оплата выполнена",
		"order_uuid", in.OrderUUID,
		"transaction_uuid", transactionUUID,
		"payment_method", in.PaymentMethod)

	return uuid.New(), nil
}
