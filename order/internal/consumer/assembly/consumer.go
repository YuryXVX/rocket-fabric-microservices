package assembly

import (
	"context"
	"log/slog"
)

type service struct {
	assemblyConsumer Consumer
	orderService     Service
}

func NewService(consumer Consumer, orderService Service) *service {
	return &service{
		orderService:     orderService,
		assemblyConsumer: consumer,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	slog.InfoContext(ctx, "запуск потребителя AssemblyService")

	return s.assemblyConsumer.Consume(ctx, s.AssemblyHandler)
}
