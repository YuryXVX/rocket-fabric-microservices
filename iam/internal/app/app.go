package app

import (
	"context"
	"iam/internal/config"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"plaform/pkg/closer"
	"plaform/pkg/grpc/health"
	v1Auth "shared/pkg/proto/auth/v1"
	v1User "shared/pkg/proto/user/v1"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	// gRPC keepalive параметры
	grpcMaxConnectionIdle     = 15 * time.Minute // Закрыть idle-соединения (нет активных RPC)
	grpcMaxConnectionAge      = 30 * time.Minute // Принудительная ротация для балансировки
	grpcMaxConnectionAgeGrace = 5 * time.Second  // Время на завершение активных RPC
	grpcKeepaliveTime         = 5 * time.Minute  // Интервал ping'ов для обнаружения мёртвых соединений
	grpcKeepaliveTimeout      = 1 * time.Second  // Таймаут ожидания pong
	grpcMinPingInterval       = 5 * time.Minute  // Минимальный интервал ping'ов от клиента (защита от DoS)
	shutdownTimeout           = 5 * time.Second
)

type App struct {
	grpcServer  *grpc.Server
	diContainer diContainer
	listener    net.Listener
}

func New(ctx context.Context) *App {
	a := &App{}

	a.initDeps(ctx)

	return a
}

func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		errCh <- a.runGRPCServer()
	}()

	var runErr error

	select {
	case runErr = <-errCh:
	case <-ctx.Done():
		slog.Info("получен сигнал завершения, начинаем graceful shutdown")
	}

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)

	defer shutdownCancel()

	if err := closer.CloseAll(shutdownCtx); err != nil {
		slog.Error("ошибка при завершении работы", "error", err)
		if runErr == nil {
			runErr = err
		}
	}

	return runErr
}

func (a *App) initDeps(ctx context.Context) {
	deps := []func(context.Context){
		a.initDI,
		a.initListener,
		a.initGRPCServer,
	}

	for _, f := range deps {
		f(ctx)
	}
}

func (a *App) initGRPCServer(ctx context.Context) {
	a.grpcServer = grpc.NewServer(
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

	authApiV1 := a.diContainer.AuthHandlers(ctx)
	userApiV1 := a.diContainer.UserHandlers(ctx)

	closer.Add("gRPC server", func(_ context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	health.RegisterService(a.grpcServer)

	v1Auth.RegisterAuthServiceServer(a.grpcServer, authApiV1)
	v1User.RegisterUserServiceServer(a.grpcServer, userApiV1)
}

func (a *App) runGRPCServer() error {
	slog.Info("🚀 gRPC-сервер запущен", "address", config.AppConfig().Grpc.Address())

	return a.grpcServer.Serve(a.listener)
}

func (a *App) initDI(ctx context.Context) {
	a.diContainer = diContainer{}
}

func (a *App) initListener(_ context.Context) {
	listener, err := net.Listen("tcp", config.AppConfig().Grpc.Address()) //nolint:noctx // net.Listen не требует контекст, адрес из конфига
	if err != nil {
		slog.Error("не удалось создать TCP-листенер", "error", err)
		os.Exit(1)
	}

	a.listener = listener
}
