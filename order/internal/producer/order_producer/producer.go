package orderproducer

import (
	"context"
	"log/slog"
	"order/internal/model"
	"plaform/pkg/kafka"
	v1 "shared/pkg/proto/events/v1"

	"google.golang.org/protobuf/proto"
)

type service struct {
	orderProducer KafkaProducer
}

func NewProducer(producer KafkaProducer) *service {
	return &service{
		orderProducer: producer,
	}
}

func (s *service) ProduceOrderRecord(ctx context.Context, event model.OrderPaidEvent) error {
	var msg = &v1.OrderPaid{
		EventUuid: event.EventUuid,
		OrderUuid: event.OrderUuid,
		UserUuid:  event.UserUuid,
	}

	payload, err := proto.Marshal(msg)

	if err != nil {
		slog.ErrorContext(ctx, "не удалось сериализовать OrderRecord", "error", err)
		return err
	}

	return s.orderProducer.Send(ctx, &kafka.Message{
		Key:   []byte(event.EventUuid),
		Value: payload,
	})
}
