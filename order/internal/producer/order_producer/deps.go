package orderproducer

import (
	"context"
	"plaform/pkg/kafka"
)

type KafkaProducer interface {
	Send(ctx context.Context, msg *kafka.Message) error
}
