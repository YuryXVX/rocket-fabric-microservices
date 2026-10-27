package app

import (
	"context"
	"log/slog"
	apiV1 "order/internal/api/order/v1"
	"order/internal/client/grpc/inventory"
	"order/internal/client/grpc/payment"
	"order/internal/consumer/assembly"
	orderproducer "order/internal/producer/order_producer"
	repository "order/internal/repository/order"
	service "order/internal/service/order"
	"os"
	"plaform/pkg/closer"
	"plaform/pkg/di"
	"plaform/pkg/kafka/consumer"
	kafka_middleware "plaform/pkg/kafka/middleware"
	inventoryV1 "shared/pkg/proto/inventory/v1"
	paymentV1 "shared/pkg/proto/payment/v1"

	"github.com/IBM/sarama"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"

	orderV1 "shared/pkg/openapi/order/v1"

	wrappedKafkaConsumer "plaform/pkg/kafka/consumer"
	"plaform/pkg/kafka/producer"
)

const (
	inventoryAddress = "127.0.0.1:50051"
	paymentAddress   = "127.0.0.1:50050"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type diContainer struct {
	// grpc clients
	inventoryGrpcClient di.Value[service.InventoryClientGrpc]
	paymentGrpcClient   di.Value[service.PaymentClientGrpc]

	// обертки для kafka
	producer di.Value[*producer.Producer]
	consumer di.Value[*consumer.Consumer]

	syncProducer  di.Value[sarama.SyncProducer]
	consumerGroup di.Value[sarama.ConsumerGroup]

	// tx manager
	txManager di.Value[repository.TxManager]
	// pool
	pgPool di.Value[*pgxpool.Pool]
	// repository
	orderRepository di.Value[service.OrderRepository]
	// service
	orderService di.Value[apiV1.OrderService]

	api di.Value[orderV1.Handler]

	// kafka services
	assemblyConsumerSvc  di.Value[ConsumerService]
	orderPaidProducerSvc di.Value[service.OrderPaidProducer]
}

func (di *diContainer) Pool(ctx context.Context) *pgxpool.Pool {
	return di.pgPool.Get(ctx, func(ctx context.Context) *pgxpool.Pool {
		dbURI := os.Getenv("DB_URI")

		if dbURI == "" {
			slog.Error("переменная окружения DB_URI не установлена")
			return nil
		}

		pool, err := pgxpool.New(ctx, dbURI)

		if err != nil {
			slog.Error("ошибка подключения к БД", "error", err)
			return nil
		}

		closer.Add("PostgreSQL pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		return pool
	})
}

func createConnection(target string) (*grpc.ClientConn, error) {
	return grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}

func (di *diContainer) Tx(ctx context.Context) repository.TxManager {
	return di.txManager.Get(ctx, func(ctx context.Context) repository.TxManager {
		tx, err := manager.New(trmpgx.NewDefaultFactory(di.Pool(ctx)))

		if err != nil {
			slog.Error("tx manager", "error", err)
			return nil
		}

		return tx
	})
}

func (di *diContainer) OrderRepository(ctx context.Context) service.OrderRepository {
	return di.orderRepository.Get(ctx, func(ctx context.Context) service.OrderRepository {
		return repository.New(di.Pool(ctx))
	})
}

func (di *diContainer) InventoryGRPCClient(ctx context.Context) service.InventoryClientGrpc {
	return di.inventoryGrpcClient.Get(ctx, func(ctx context.Context) service.InventoryClientGrpc {
		connection, err := createConnection(inventoryAddress)

		if err != nil {
			slog.Error("ошибки при создании коннекшена", "error", err)
		}

		closer.Add("close inventory grpc connection", func(ctx context.Context) error {
			connection.Close()
			return nil
		})

		client := inventoryV1.NewPartServiceClient(connection)

		return inventory.New(client)
	})
}

func (di *diContainer) PaymentGRPCClient(ctx context.Context) service.PaymentClientGrpc {
	return di.paymentGrpcClient.Get(ctx, func(ctx context.Context) service.PaymentClientGrpc {
		connection, err := createConnection(paymentAddress)

		if err != nil {
			slog.Error("ошибки при создании коннекшена", "error", err)
		}

		closer.Add("close payment grpc connection", func(ctx context.Context) error {
			connection.Close()
			return nil
		})

		client := paymentV1.NewBillingServiceClient(connection)

		return payment.New(client)
	})
}

func (di *diContainer) OrderService(ctx context.Context) apiV1.OrderService {
	return di.orderService.Get(ctx, func(ctx context.Context) apiV1.OrderService {
		return service.New(
			di.InventoryGRPCClient(ctx),
			di.PaymentGRPCClient(ctx),
			di.OrderRepository(ctx),
			di.Tx(ctx),
			di.OrderPaidProduce(ctx),
		)
	})
}

func (di *diContainer) OrderHandlers(ctx context.Context) orderV1.Handler {
	return di.api.Get(ctx, func(ctx context.Context) orderV1.Handler {
		return apiV1.New(di.OrderService(ctx))
	})
}

func (di *diContainer) OrderPaidProduce(ctx context.Context) service.OrderPaidProducer {
	return di.orderPaidProducerSvc.Get(ctx, func(ctx context.Context) service.OrderPaidProducer {
		return orderproducer.NewProducer(di.Producer(ctx))
	})
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
		return producer.NewProducer(d.SyncProducer(ctx), "order.paid")
	})
}

func (d *diContainer) Consumer(ctx context.Context) *wrappedKafkaConsumer.Consumer {
	return d.consumer.Get(ctx, func(ctx context.Context) *wrappedKafkaConsumer.Consumer {
		return wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(ctx),
			[]string{
				"assembly.ship-assembled",
			},
			wrappedKafkaConsumer.WithMiddlewares(
				kafka_middleware.ConsumerLogging(),
			),
		)
	})
}

func (d *diContainer) AssemblyServiceProducer(ctx context.Context) ConsumerService {
	return d.assemblyConsumerSvc.Get(ctx, func(ctx context.Context) ConsumerService {
		return assembly.NewService(d.Consumer(ctx), d.OrderService(ctx))
	})
}

func (d *diContainer) ConsumerGroup(ctx context.Context) sarama.ConsumerGroup {
	return d.consumerGroup.Get(ctx, func(ctx context.Context) sarama.ConsumerGroup {
		consumerGroup, err := sarama.NewConsumerGroup(
			[]string{"localhost:9092"},
			"1",
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

func SaramaConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V4_0_0_0
	cfg.Producer.Return.Successes = true

	cfg.Metadata.AllowAutoTopicCreation = true

	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	return cfg
}
