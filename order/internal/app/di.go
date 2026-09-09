package app

import (
	"context"
	"log/slog"
	apiV1 "order/internal/api/order/v1"
	"order/internal/client/grpc/inventory"
	"order/internal/client/grpc/payment"
	repository "order/internal/repository/order"
	service "order/internal/service/order"
	"os"
	"plaform/pkg/closer"
	"plaform/pkg/di"
	inventoryV1 "shared/pkg/proto/inventory/v1"
	paymentV1 "shared/pkg/proto/payment/v1"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"

	orderV1 "shared/pkg/openapi/order/v1"
)

const (
	inventoryAddress = "localhost:50051"
	paymentAddress   = "localhost:50050"
)

type diContainer struct {
	// grpc clients
	inventoryGrpcClient di.Value[service.InventoryClientGrpc]
	paymentGrpcClient   di.Value[service.PaymentClientGrpc]

	// tx manager
	txManager di.Value[repository.TxManager]
	// pool
	pgPool di.Value[*pgxpool.Pool]
	// repository
	orderRepository di.Value[service.OrderRepository]
	// service
	orderService di.Value[apiV1.OrderService]

	api di.Value[orderV1.Handler]
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
		return repository.New(di.Pool(ctx), di.Tx(ctx))
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
		)
	})
}

func (di *diContainer) OrderHandlers(ctx context.Context) orderV1.Handler {
	return di.api.Get(ctx, func(ctx context.Context) orderV1.Handler {
		return apiV1.New(di.OrderService(ctx))
	})
}
