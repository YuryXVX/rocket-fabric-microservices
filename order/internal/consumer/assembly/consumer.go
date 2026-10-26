package assembly

import (
	"context"
	"log/slog"
)

type service struct {
	assemblyConsumer Consumer
}

func NewService(consumer Consumer) *service {
	return &service{
		assemblyConsumer: consumer,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	slog.InfoContext(ctx, "запуск потребителя UFORecorded")

	return s.assemblyConsumer.Consume(ctx, s.AssemblyHandler)
}
