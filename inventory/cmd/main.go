package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "inventory/internal/api/inventory/v1"
	repository "inventory/internal/repository/inventory"
	service "inventory/internal/service/inventory"
	v1 "shared/pkg/proto/inventory/v1"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

var (
	// адрес сервера
	grpcAddress = "localhost:50051"

	// gRPC keepalive параметры
	grpcMaxConnectionIdle     = 15 * time.Minute // Закрыть idle-соединения (нет активных RPC)
	grpcMaxConnectionAge      = 30 * time.Minute // Принудительная ротация для балансировки
	grpcMaxConnectionAgeGrace = 5 * time.Second  // Время на завершение активных RPC
	grpcKeepaliveTime         = 5 * time.Minute  // Интервал ping'ов для обнаружения мёртвых соединений
	grpcKeepaliveTimeout      = 1 * time.Second  // Таймаут ожидания pong
	grpcMinPingInterval       = 5 * time.Minute  // Минимальный интервал ping'ов от клиента (защита от DoS)
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	lis, err := net.Listen("tcp", grpcAddress)

	if err != nil {
		slog.Error("ошибка запуска слушателя", "error", err)

		return
	}

	err = godotenv.Load("../../inventory.env")

	if err != nil {
		slog.Error("ошибка загрузки переменных окружения из inventory.env", "error", err)

		return
	}

	dbURI := os.Getenv("DB_URI")

	if dbURI == "" {
		slog.Error("переменная окружения DB_URI не установлена")
		return
	}

	pool, err := pgxpool.New(ctx, dbURI)

	if err != nil {
		slog.Error("ошибка подключения к БД", "error", err)
		return
	}
	defer pool.Close()

	s := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     grpcMaxConnectionIdle,
			MaxConnectionAge:      grpcMaxConnectionAge,
			MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
			Time:                  grpcKeepaliveTime,
			Timeout:               grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcMinPingInterval,
			PermitWithoutStream: true, // Разрешить "тёплые" соединения без активных RPC
		}),
	)

	repository := repository.NewRepository(pool)
	service := service.NewService(repository)
	api := api.New(service)

	v1.RegisterPartServiceServer(s, api)

	reflection.Register(s)

	go func() {
		slog.Info("🚀 gRPC сервер запущен", "address", grpcAddress)

		if err := s.Serve(lis); err != nil {
			slog.Error("ошибка запуска сервера", "error", err)
			cancel()
		}
	}()

	<-ctx.Done()
	slog.Info("🛑 остановка gRPC сервера")
	s.GracefulStop()

	slog.Info("✅ сервер остановлен")
}
