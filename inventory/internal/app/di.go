package app

import (
	"context"
	apiV1 "inventory/internal/api/inventory/v1"
	repositoryInventory "inventory/internal/repository/inventory"
	application "inventory/internal/service/application/part"
	"inventory/internal/service/domain"
	"log/slog"
	"os"
	"plaform/pkg/closer"
	"plaform/pkg/di"
	v1 "shared/pkg/proto/inventory/v1"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/jackc/pgx/v5/pgxpool"
)

type diContainer struct {
	// poll
	pgPool di.Value[*pgxpool.Pool]
	// repository
	inventoryRepository di.Value[application.InventoryRepository]
	// service
	applicationService di.Value[apiV1.ApplicationService]
	// handlers
	inventoryHandlers di.Value[v1.PartServiceServer]
	// domain service checker part
	compatibilityPartChecker di.Value[application.CompatibilityChecker]
	// tx manager
	txManager di.Value[application.TxManager]
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

		if err := pool.Ping(ctx); err != nil {
			slog.Error("База данных недоступна: %v", err)
		}

		closer.Add("PostgreSQL pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		return pool
	})
}

func (di *diContainer) InventoryRepository(ctx context.Context) application.InventoryRepository {
	return di.inventoryRepository.Get(ctx, func(ctx context.Context) application.InventoryRepository {
		return repositoryInventory.NewRepository(di.PgPoll(ctx))
	})
}

func (di *diContainer) ApplicationService(ctx context.Context) apiV1.ApplicationService {
	return di.applicationService.Get(ctx, func(ctx context.Context) apiV1.ApplicationService {
		return application.NewService(
			di.InventoryRepository(ctx),
			di.CompatibilityChecker(ctx),
			di.TxManager(ctx),
		)
	})
}

func (di *diContainer) InventoryHandlers(ctx context.Context) v1.PartServiceServer {
	return di.inventoryHandlers.Get(ctx, func(ctx context.Context) v1.PartServiceServer {
		return apiV1.New(di.ApplicationService(ctx))
	})
}

func (di *diContainer) CompatibilityChecker(ctx context.Context) application.CompatibilityChecker {
	return di.compatibilityPartChecker.Get(ctx, func(ctx context.Context) application.CompatibilityChecker {
		return domain.New()
	})
}

func (di *diContainer) TxManager(ctx context.Context) application.TxManager {
	return di.txManager.Get(ctx, func(ctx context.Context) application.TxManager {
		tx, err := manager.New(trmpgx.NewDefaultFactory(di.PgPoll(ctx)))

		if err != nil {
			slog.Error("tx manager", "error", err)
			return nil
		}

		return tx
	})
}
