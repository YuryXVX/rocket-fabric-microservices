package kafka_middleware

import (
	"context"
	"log/slog"
	"plaform/pkg/kafka"
)

// ConsumerLogging — middleware для логирования входящих сообщений
func ConsumerLogging() kafka.Middleware {
	return func(next kafka.MessageHandler) kafka.MessageHandler {
		return func(ctx context.Context, msg kafka.Message) error {
			slog.InfoContext(ctx, "получено сообщение Kafka", "topic", msg.Topic)
			return next(ctx, msg)
		}
	}
}
