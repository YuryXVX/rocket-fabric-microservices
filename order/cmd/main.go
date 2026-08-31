package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	customMiddleware "order/internal/api/middleware"
	apiV1 "order/internal/api/order/v1"
	"order/internal/client/grpc/inventory"
	"order/internal/client/grpc/payment"
	"order/internal/repository/order"
	service "order/internal/service/order"
	orderV1 "shared/pkg/openapi/order/v1"
	inventoryV1 "shared/pkg/proto/inventory/v1"
	paymentV1 "shared/pkg/proto/payment/v1"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer cancel()

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

	err = godotenv.Load("../order.env")

	if err != nil {
		slog.Error("ошибка загрузки переменных окружения из order.env", "error", err)

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

	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))

	serviceInventory := inventoryV1.NewPartServiceClient(inventoryConn)
	servicePayment := paymentV1.NewBillingServiceClient(paymentConn)

	inventoryClient := inventory.New(serviceInventory)
	paymentClient := payment.New(servicePayment)

	orderRepository := order.New(pool, txManager)

	service := service.New(
		inventoryClient,
		paymentClient,
		orderRepository,
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
}
