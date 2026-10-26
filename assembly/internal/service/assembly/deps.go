package assembly

import (
	"assemply/internal/model"
	"context"
	"plaform/pkg/kafka"
	"time"
)

type OrderPaidConsumer interface {
	OrderPaidHandler(ctx context.Context, msg kafka.Message) (model.OrderPaidEvent, error)
	OrderPaidRunConsumer(ctx context.Context, handler kafka.MessageHandler) error
}

type AssemblyProducer interface {
	AssemblyEventProduced(ctx context.Context, event model.OrderPaidEvent, start time.Time) error
}
