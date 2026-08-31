package order

import (
	"context"
	"order/internal/model"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) error
	Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	Update(ctx context.Context, order model.Order) error
}

type PartRepository interface {
	GetPart(ctx context.Context, orderUUID uuid.UUID) ([]model.OrderItem, error)
	SavePart(ctx context.Context, order model.Order) error
}

type InventoryClientGrpc interface {
	ListParts(ctx context.Context, in []uuid.UUID) ([]model.OrderItem, error)
}

type PaymentClientGrpc interface {
	PayOrder(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (uuid.UUID, error)
}
