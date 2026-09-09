package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"plaform/pkg/logger"
	orderV1 "shared/pkg/openapi/order/v1"
	"syscall"
	"time"

	customMiddleware "order/internal/api/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	httpPort = "8080"

	// // Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	shutdownTimeout   = 10 * time.Second
	middlewareTimeout = 100 * time.Second
)

type App struct {
	diContainer *diContainer
}

func New(cxt context.Context) *App {
	a := &App{}

	a.initDeps(cxt)

	return a
}

func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	r, err := a.setupRouter(a.diContainer.OrderHandlers(ctx))

	if err != nil {
		slog.Error("ошибка при создании роутера", "error", err)
		return err
	}

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	go func() {
		slog.Info("🚀 старт сервера", "port", httpPort)
		if serveErr := server.ListenAndServe(); serveErr != nil && serveErr != http.ErrServerClosed {
			slog.Error("❌ ошибка запуска сервера", "error", serveErr)
			cancel()
		}
	}()

	<-ctx.Done()
	slog.Info("🛑 завершение работы сервера...")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), time.Hour)
	defer cancelShutdown()

	if shutdownErr := server.Shutdown(shutdownCtx); shutdownErr != nil {
		if errors.Is(shutdownErr, context.DeadlineExceeded) {
			slog.Error("❌ время остановки сервера истекло", "timeout", time.Hour)
		} else {
			slog.Error("❌ ошибка при остановке сервера", "error", shutdownErr)
		}
	}

	slog.Info("✅ сервер остановлен")

	return nil
}

func (a *App) initDeps(ctx context.Context) {
	inits := []func(context.Context){
		a.initDI,
		a.initLogger,
	}

	for _, f := range inits {
		f(ctx)
	}
}

func (a *App) initDI(_ context.Context) {
	a.diContainer = &diContainer{}
}

func (a *App) initLogger(ctx context.Context) {
	logger.Init("info")
}

func (a *App) setupRouter(h orderV1.Handler) (chi.Router, error) {
	orderServer, err := orderV1.NewServer(h)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания сервера OpenAPI: %w", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(customMiddleware.RequestLogger)
	r.Use(middleware.Timeout(middlewareTimeout))

	r.Handle("/api/*", orderServer)

	return r, nil
}
