package assembly

import (
	"context"
	"plaform/pkg/kafka"

	"github.com/google/uuid"
)

type Consumer interface {
	Consume(ctx context.Context, handler kafka.MessageHandler) error
}

type Service interface {
	CommitOrder(ctx context.Context, orderUUID uuid.UUID) error
}
