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

	customMiddleware "order/internal/api/middleware"
	apiV1 "order/internal/api/order/v1"
	"order/internal/client/grpc/inventory"
	"order/internal/client/grpc/payment"
	"order/internal/repository/order"
	"order/internal/repository/part"
	service "order/internal/service/order"
	orderV1 "shared/pkg/openapi/order/v1"
	inventoryV1 "shared/pkg/proto/inventory/v1"
	paymentV1 "shared/pkg/proto/payment/v1"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	inventoryAddress = "localhost:50051"
	paymentAddress   = "localhost:50050"
)

func setupRouter(h orderV1.Handler) (chi.Router, error) {
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

func main() {
	paymentConn, err := grpc.NewClient(paymentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		slog.Error("произошла ошибка при создании клиента payment", "error", err)
		return
	}

	defer func() {
		if err := paymentConn.Close(); err != nil {
			slog.Error("ошибка закрытия соединения", "error", err)
		}
	}()

	inventoryConn, err := grpc.NewClient(inventoryAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		slog.Error("произошла ошибка при создании клиента inventory", "error", err)
		return
	}

	defer func() {
		if err := inventoryConn.Close(); err != nil {
			slog.Error("ошибка закрытия соединения", "error", err)
		}
	}()

	serviceInventory := inventoryV1.NewPartServiceClient(inventoryConn)
	servicePayment := paymentV1.NewBillingServiceClient(paymentConn)

	inventoryClient := inventory.New(serviceInventory)
	paymentClient := payment.New(servicePayment)
	orderRepository := order.New()
	partRepository := part.New()

	service := service.New(
		inventoryClient,
		paymentClient,
		orderRepository,
		partRepository,
	)

	handler := apiV1.New(service)

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
}
