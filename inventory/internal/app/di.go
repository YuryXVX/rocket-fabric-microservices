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
	v1 "shared/pkg/proto/inventory/v1"

	"github.com/jackc/pgx/v5/pgxpool"
)

type diContainer struct {
	// poll
	pgPool *pgxpool.Pool
	// repository
	inventoryRepository serviceInventory.InventoryRepository
	// service
	inventoryService apiV1.ServiceInventory
	// handlers
	inventoryHandlers v1.PartServiceServer
}

func (di *diContainer) PgPoll(ctx context.Context) *pgxpool.Pool {
	if di.pgPool == nil {
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

		di.pgPool = pool

	}

	return di.pgPool
}

func (di *diContainer) InventoryRepository(ctx context.Context) inventory.InventoryRepository {
	if di.inventoryRepository == nil {
		di.inventoryRepository = repositoryInventory.NewRepository(di.PgPoll(ctx))
	}

	return di.inventoryRepository
}

func (di *diContainer) InventoryService(ctx context.Context) apiV1.ServiceInventory {
	if di.inventoryService == nil {
		di.inventoryService = serviceInventory.NewService(di.InventoryRepository(ctx))
	}

	return di.inventoryService
}

func (di *diContainer) InventoryHandlers(ctx context.Context) v1.PartServiceServer {
	if di.inventoryHandlers == nil {
		di.inventoryHandlers = apiV1.New(di.InventoryService(ctx))
	}

	return di.inventoryHandlers
}
