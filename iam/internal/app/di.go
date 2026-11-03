package app

import (
	"context"
	apiSession "iam/internal/api/session"
	apiUser "iam/internal/api/user"
	"iam/internal/config"
	"iam/internal/infrastructure/crypto"
	sessionRepository "iam/internal/repository/session"
	userRepository "iam/internal/repository/user"
	iamService "iam/internal/service/iam"
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
	apiUser.UserService
	apiSession.SessionService
}

type diContainer struct {
	pgPool      di.Value[*pgxpool.Pool]
	redisClient di.Value[*redis.Client]

	authHandlers di.Value[v1Auth.AuthServiceServer]
	userHandlers di.Value[v1User.UserServiceServer]

	iamService di.Value[IamService]

	userRepository    di.Value[iamService.UserRepository]
	sessionRepository di.Value[iamService.SessionRepository]

	hasher di.Value[iamService.Hasher]
}

func (di *diContainer) PgPoll(ctx context.Context) *pgxpool.Pool {
	return di.pgPool.Get(ctx, func(ctx context.Context) *pgxpool.Pool {
		pool, err := pgxpool.New(ctx, config.AppConfig().PG.DSN())

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
			Addr:            config.AppConfig().Redis.Address(),
			DialTimeout:     config.AppConfig().Redis.ConnectionTimeout,
			ReadTimeout:     config.AppConfig().Redis.ConnectionTimeout,
			WriteTimeout:    config.AppConfig().Redis.ConnectionTimeout,
			MaxIdleConns:    config.AppConfig().Redis.MaxIdle,
			ConnMaxIdleTime: config.AppConfig().Redis.IdleTimeout,
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
		return apiSession.New(di.IamService(ctx))
	})
}

func (di *diContainer) UserHandlers(ctx context.Context) v1User.UserServiceServer {
	return di.userHandlers.Get(ctx, func(ctx context.Context) v1User.UserServiceServer {
		return apiUser.New(di.IamService(ctx))
	})
}

func (di *diContainer) IamService(ctx context.Context) IamService {
	return di.iamService.Get(ctx, func(ctx context.Context) IamService {
		return iamService.New(
			di.UserRepository(ctx),
			di.SessionRepository(ctx),
			di.Hasher(ctx),
		)
	})
}

func (di *diContainer) UserRepository(ctx context.Context) iamService.UserRepository {
	return di.userRepository.Get(ctx, func(ctx context.Context) iamService.UserRepository {
		return userRepository.New(di.PgPoll(ctx))
	})
}

func (di *diContainer) SessionRepository(ctx context.Context) iamService.SessionRepository {
	return di.sessionRepository.Get(ctx, func(ctx context.Context) iamService.SessionRepository {
		return sessionRepository.New(di.RedisClient(ctx))
	})
}

func (di *diContainer) Hasher(ctx context.Context) iamService.Hasher {
	return di.hasher.Get(ctx, func(ctx context.Context) iamService.Hasher {
		return crypto.NewBcryptHasher()
	})
}
