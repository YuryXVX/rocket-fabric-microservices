package assebmlyproducer

import (
	"assemply/internal/model"
	"context"
	"log/slog"
	"plaform/pkg/kafka"
	"plaform/pkg/kafka/producer"
	v1 "shared/pkg/proto/events/v1"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service struct {
	producer *producer.Producer
}

func NewService(producer *producer.Producer) *service {
	return &service{
		producer: producer,
	}
}

func (s *service) AssemblyEventProduced(ctx context.Context, event model.OrderPaidEvent, startBuild time.Time) error {
	msg := &v1.ShipAssembled{
		EventUuid:    uuid.NewString(),
		OrderUuid:    event.OrderUuid,
		UserUuid:     event.UserUuid,
		BuildTimeSec: int64(time.Since(startBuild).Seconds()),
		AssembledAt:  timestamppb.Now(),
	}

	payload, err := proto.Marshal(msg)

	if err != nil {
		slog.ErrorContext(ctx, "не удалось сериализовать ShipAssembled", "error", err)
		return err
	}

	return s.producer.Send(ctx, &kafka.Message{
		Key:   []byte(msg.EventUuid),
		Value: payload,
	})
}
