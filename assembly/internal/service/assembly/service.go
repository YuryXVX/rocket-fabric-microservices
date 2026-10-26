package assembly

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"plaform/pkg/kafka"
	"time"
)

type service struct {
	orderPaidConsumer OrderPaidConsumer
	assemblyProducer  AssemblyProducer
}

func NewService(orderPaidConsumer OrderPaidConsumer, assemblyProducer AssemblyProducer) *service {
	return &service{
		orderPaidConsumer: orderPaidConsumer,
		assemblyProducer:  assemblyProducer,
	}
}

func (s *service) Process(ctx context.Context) error {
	slog.InfoContext(ctx, "запуск потребителя Order.Paid")

	return s.orderPaidConsumer.OrderPaidRunConsumer(ctx, s.ProcessHandler)
}

func (s *service) ProcessHandler(ctx context.Context, msg kafka.Message) error {
	start := time.Now()

	paidEvent, err := s.orderPaidConsumer.OrderPaidHandler(ctx, msg)

	if err != nil {
		return fmt.Errorf("ошибка при чтении сообщения Order.Paid: %w", err)
	}

	slog.InfoContext(ctx, "запуск производителя Assembly.Event")

	randDuration := rand.IntN(15)

	time.Sleep(time.Duration(randDuration) * time.Second)

	if err = s.assemblyProducer.AssemblyEventProduced(ctx, paidEvent, start); err != nil {
		return fmt.Errorf("ошибка при запуске производителя Assembly.Event %w", err)
	}

	return nil
}
