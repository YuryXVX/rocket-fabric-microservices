package order

import (
	"context"
	"order/internal/model"

	"github.com/google/uuid"
)

func (s *service) Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	return uuid.New(), nil
}
