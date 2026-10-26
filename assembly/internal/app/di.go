package app

import (
	orderpaidconsumer "assemply/internal/consumer/order_paid"
	assebmlyproducer "assemply/internal/producer/assembly"
	"assemply/internal/service/assembly"
	"context"
	"log/slog"
	"os"
	"plaform/pkg/closer"
	"plaform/pkg/di"
	"plaform/pkg/kafka/consumer"
	"plaform/pkg/kafka/producer"

	"github.com/IBM/sarama"
)

type AssemblyService interface {
	Process(ctx context.Context) error
}

type diContainer struct {
	// обертки для kafka
	producer      di.Value[*producer.Producer]
	consumer      di.Value[*consumer.Consumer]
	syncProducer  di.Value[sarama.SyncProducer]
	consumerGroup di.Value[sarama.ConsumerGroup]

	orderPaidConsumer di.Value[assembly.OrderPaidConsumer]
	assemblyProducer  di.Value[assembly.AssemblyProducer]

	assemblyService di.Value[AssemblyService]
}

func (d *diContainer) SyncProducer(ctx context.Context) sarama.SyncProducer {
	return d.syncProducer.Get(ctx, func(ctx context.Context) sarama.SyncProducer {
		config := SaramaConfig()

		p, err := sarama.NewSyncProducer(
			[]string{"localhost:9092"},
			config,
		)

		if err != nil {
			slog.Error("не удалось создать sync producer", "error", err)
			os.Exit(1)
		}

		closer.Add("Kafka sync producer", func(_ context.Context) error {
			return p.Close()
		})

		return p
	})
}

func (d *diContainer) Producer(ctx context.Context) *producer.Producer {
	return d.producer.Get(ctx, func(ctx context.Context) *producer.Producer {
		return producer.NewProducer(d.SyncProducer(ctx), "assembly.ship-assembled")
	})
}

func (d *diContainer) Consumer(ctx context.Context) *consumer.Consumer {
	return d.consumer.Get(ctx, func(ctx context.Context) *consumer.Consumer {
		return consumer.NewConsumer(
			d.ConsumerGroup(ctx),
			[]string{"order.paid"},
		)
	})
}

func (d *diContainer) ConsumerGroup(ctx context.Context) sarama.ConsumerGroup {
	return d.consumerGroup.Get(ctx, func(ctx context.Context) sarama.ConsumerGroup {
		consumerGroup, err := sarama.NewConsumerGroup(
			[]string{"localhost:9092"},
			"2",
			SaramaConfig(),
		)
		if err != nil {
			slog.Error("не удалось создать consumer group", "error", err)
			os.Exit(1)
		}

		closer.Add("Kafka consumer group", func(_ context.Context) error {
			return consumerGroup.Close()
		})

		return consumerGroup
	})
}

func (d *diContainer) AssemblyService(ctx context.Context) AssemblyService {
	return d.assemblyService.Get(ctx, func(ctx context.Context) AssemblyService {
		return assembly.NewService(d.OrderPaidConsumer(ctx), d.AssemblyProducer(ctx))
	})
}

func (d *diContainer) OrderPaidConsumer(ctx context.Context) assembly.OrderPaidConsumer {
	return d.orderPaidConsumer.Get(ctx, func(ctx context.Context) assembly.OrderPaidConsumer {
		return orderpaidconsumer.NewService(d.Consumer(ctx))
	})
}

func (d *diContainer) AssemblyProducer(ctx context.Context) assembly.AssemblyProducer {
	return d.assemblyProducer.Get(ctx, func(ctx context.Context) assembly.AssemblyProducer {
		return assebmlyproducer.NewService(d.Producer(ctx))
	})
}

func SaramaConfig() *sarama.Config {
	cfg := sarama.NewConfig()

	cfg.Producer.Return.Successes = true
	cfg.Version = sarama.V4_0_0_0
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	return cfg
}
