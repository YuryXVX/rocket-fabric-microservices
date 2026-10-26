package orderpaidconsumer

import (
	"assemply/internal/model"
	"context"
	"fmt"
	"plaform/pkg/kafka"
	"plaform/pkg/kafka/consumer"
	v1 "shared/pkg/proto/events/v1"

	"google.golang.org/protobuf/proto"
)

type service struct {
	consumer *consumer.Consumer
}

func NewService(c *consumer.Consumer) *service {
	return &service{
		consumer: c,
	}
}

func decodeOrderPaidMessage(mgs kafka.Message) (model.OrderPaidEvent, error) {
	var pb v1.OrderPaid
	if err := proto.Unmarshal(mgs.Value, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("ошибка при Unmarshal order paid %w", err)
	}

	return model.OrderPaidEvent{
		UserUuid:  pb.UserUuid,
		EventUuid: pb.EventUuid,
		OrderUuid: pb.OrderUuid,
	}, nil
}

func (c *service) OrderPaidHandler(ctx context.Context, mgs kafka.Message) (model.OrderPaidEvent, error) {
	event, err := decodeOrderPaidMessage(mgs)

	if err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("ошибка при декодировании сообщения %w", err)
	}

	fmt.Printf("[mgs] - %+v", event)

	return event, nil
}

func (c *service) OrderPaidRunConsumer(ctx context.Context, handler kafka.MessageHandler) error {
	return c.consumer.Consume(ctx, handler)
}
