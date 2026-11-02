package app

import (
	"context"
	apiAuthV1 "iam/internal/api/iam"
	api_session "iam/internal/api/iam"
	api_user "iam/internal/api/user"
	iam_service "iam/internal/service/iam"
	"log/slog"
	"os"
	"plaform/pkg/closer"
	"plaform/pkg/di"
	platform_redis "plaform/pkg/redis"
	v1Auth "shared/pkg/proto/auth/v1"
	v1User "shared/pkg/proto/user/v1"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type IamService interface {
	api_user.UserService
	api_session.SessionService
}

type diContainer struct {
	pgPool      di.Value[*pgxpool.Pool]
	redisClient di.Value[*redis.Client]

	authHandlers di.Value[v1Auth.AuthServiceServer]
	userHandlers di.Value[v1User.UserServiceServer]

	iamService di.Value[IamService]

	userRepository    di.Value[iam_service.UserRepository]
	sessionRepository di.Value[iam_service.SessionRepository]
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
			slog.Error("База данных недоступна: %+w", err)
		}

		closer.Add("PostgreSQL pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		return pool
	})
}

func (di *diContainer) RedisClient(ctx context.Context) *redis.Client {
	return di.redisClient.Get(ctx, func(ctx context.Context) *redis.Client {
		rdb, err := platform_redis.NewClient(&redis.Options{
			// Addr:            "",
			// DialTimeout:     "",
			// ReadTimeout:     config.AppConfig().Redis.ConnectionTimeout,
			// WriteTimeout:    config.AppConfig().Redis.ConnectionTimeout,
			// MaxIdleConns:    config.AppConfig().Redis.MaxIdle,
			// ConnMaxIdleTime: config.AppConfig().Redis.IdleTimeout,
		}, slog.Default())
		if err != nil {
			slog.Error("не удалось создать Redis клиент", "error", err)
			os.Exit(1)
		}

		closer.Add("Redis", func(_ context.Context) error {
			return rdb.Close()
		})

		return rdb
	})
}

func (di *diContainer) AuthHandlers(ctx context.Context) v1Auth.AuthServiceServer {
	return di.authHandlers.Get(ctx, func(ctx context.Context) v1Auth.AuthServiceServer {
		return apiAuthV1.New(di.IamService(ctx))
	})
}

func (di *diContainer) UserHandlers(ctx context.Context) v1User.UserServiceServer {
	return di.userHandlers.Get(ctx, func(ctx context.Context) v1User.UserServiceServer {
		return api_user.New(di.IamService(ctx))
	})
}

func (di *diContainer) IamService(ctx context.Context) IamService {
	return di.iamService.Get(ctx, func(ctx context.Context) IamService {
		return iam_service.New()
	})
}
