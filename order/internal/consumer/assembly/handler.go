package assembly

import (
	"context"
	"log/slog"
	"plaform/pkg/kafka"
)

func (s *service) AssemblyHandler(ctx context.Context, msg kafka.Message) error {
	event, err := decodeAssemblyRecorded(msg.Value)
	if err != nil {
		slog.ErrorContext(ctx, "не удалось декодировать UFORecorded", "error", err)
		return err
	}

	slog.InfoContext(
		ctx, "обработка сообщения",
		"topic", msg.Topic,
		"partition", msg.Partition,
		"offset", msg.Offset,
		"event_uuid", event.EventUuid,
		"order_uuid", event.OrderUuid,
		"user_uuid", event.UserUuid,
	)

	return nil
}
