package v1

import (
	"context"
	"order/internal/model"
	"order/internal/service/input"

	"github.com/google/uuid"
)

type OrderService interface {
	Create(ctx context.Context, in *input.CreateOrderInput) (model.Order, error)
	Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error)
	Cancel(ctx context.Context, orderUUID uuid.UUID) error
}
