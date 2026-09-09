package app

import (
	"context"
	apiV1 "inventory/internal/api/inventory/v1"
	repositoryInventory "inventory/internal/repository/inventory"
	"inventory/internal/service/inventory"
	serviceInventory "inventory/internal/service/inventory"
	"log/slog"
	"os"
	"plaform/pkg/closer"
	"plaform/pkg/di"
	v1 "shared/pkg/proto/inventory/v1"

	"github.com/jackc/pgx/v5/pgxpool"
)

type diContainer struct {
	// poll
	pgPool di.Value[*pgxpool.Pool]
	// repository
	inventoryRepository di.Value[serviceInventory.InventoryRepository]
	// service
	inventoryService di.Value[apiV1.ServiceInventory]
	// handlers
	inventoryHandlers di.Value[v1.PartServiceServer]
}

func (di *diContainer) PgPoll(ctx context.Context) *pgxpool.Pool {
	return di.pgPool.Get(ctx, func(ctx context.Context) *pgxpool.Pool {
		dbURI := os.Getenv("DB_URI")

		if dbURI == "" {
			slog.Error("переменная окружения DB_URI не установлена")
			os.Exit(1)
		}

		pool, err := pgxpool.New(ctx, dbURI)

		if err != nil {
			slog.Error("ошибка подключения к БД", "error", err)
			os.Exit(1)
		}

		err = pool.Ping(ctx)

		if err != nil {
			slog.Error("не удалось выполнить ping PostgreSQL", "error", err)
			os.Exit(1)
		}

		closer.Add("PostgreSQL pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		return pool
	})
}

func (di *diContainer) InventoryRepository(ctx context.Context) inventory.InventoryRepository {
	return di.inventoryRepository.Get(ctx, func(ctx context.Context) serviceInventory.InventoryRepository {
		return repositoryInventory.NewRepository(di.PgPoll(ctx))
	})
}

func (di *diContainer) InventoryService(ctx context.Context) apiV1.ServiceInventory {
	return di.inventoryService.Get(ctx, func(ctx context.Context) apiV1.ServiceInventory {
		return serviceInventory.NewService(di.InventoryRepository(ctx))
	})
}

func (di *diContainer) InventoryHandlers(ctx context.Context) v1.PartServiceServer {
	return di.inventoryHandlers.Get(ctx, func(ctx context.Context) v1.PartServiceServer {
		return apiV1.New(di.InventoryService(ctx))
	})
}
