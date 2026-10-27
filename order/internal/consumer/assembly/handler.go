package assembly

import (
	"context"
	"log/slog"
	"plaform/pkg/kafka"

	"github.com/google/uuid"
)

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func (s *service) AssemblyHandler(ctx context.Context, msg kafka.Message) error {
	event, err := decodeAssemblyRecorded(msg.Value)
	if err != nil {
		slog.ErrorContext(ctx, "не удалось декодировать AssemblyRecorded", "error", err)
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

	orderUUID, err := parseUUID(event.OrderUuid)

	if err != nil {
		slog.ErrorContext(ctx, "не удалось получить event.OrderUUID", "error", err)

		return err
	}

	err = s.orderService.CommitOrder(ctx, orderUUID)

	if err != nil {
		slog.ErrorContext(ctx, "ошибка при вызове orderService.CommitOrder", "error", err)
		return err
	}

	slog.InfoContext(ctx, "заказ успешно закоммичен")

	return nil
}
