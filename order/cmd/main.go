package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	customMiddleware "order/internal/api/middleware"
	v1 "order/internal/api/order/v1"
	orderv1 "shared/pkg/openapi/order/v1"
)

const (
	httpPort = "8080"

	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	shutdownTimeout   = 10 * time.Second
	middlewareTimeout = 10 * time.Second
)

func setupRouter(h orderv1.Handler) (chi.Router, error) {
	orderServer, err := orderv1.NewServer(h)
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

func main() {
	handler := v1.New()

	r, err := setupRouter(handler)
	if err != nil {
		slog.Error("ошибка при создании роутера", "error", err)
		return
	}

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атаки
		ReadTimeout:       readTimeout,       // Лимит на чтение всего запроса
		WriteTimeout:      writeTimeout,      // Лимит на запись ответа
		IdleTimeout:       idleTimeout,       // Таймаут keep-alive соединений
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		slog.Info("🚀 старт сервера", "port", httpPort)
		if serveErr := server.ListenAndServe(); serveErr != nil && serveErr != http.ErrServerClosed {
			slog.Error("❌ ошибка запуска сервера", "error", serveErr)
			cancel()
		}
	}()

	<-ctx.Done()
	slog.Info("🛑 завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelShutdown()

	if shutdownErr := server.Shutdown(shutdownCtx); shutdownErr != nil {
		if errors.Is(shutdownErr, context.DeadlineExceeded) {
			// Проверяем, является ли ошибка результатом истечения таймаута
			slog.Error("❌ время остановки сервера истекло", "timeout", shutdownTimeout)
		} else {
			slog.Error("❌ ошибка при остановке сервера", "error", shutdownErr)
		}
	}

	slog.Info("✅ сервер остановлен")
}
