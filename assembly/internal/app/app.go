package app

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"plaform/pkg/closer"
	"syscall"
	"time"
)

var (
	shutdownTimeout = 5 * time.Second
)

type App struct {
	diContainer *diContainer
}

func NewApp(ctx context.Context) *App {
	a := &App{}

	a.initDeps(ctx)

	return a
}

func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errChan <- fmt.Errorf("потребитель упал %w", err)
			return
		}

		errChan <- nil
	}()

	var runErr error
	select {
	case runErr = <-errChan:
	case <-ctx.Done():
		slog.Info("получен сигнал завершения, начинаем graceful shutdown")

	}

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

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
	inits := []func(context.Context){
		a.initDI,
	}

	for _, f := range inits {
		f(ctx)
	}
}

func (a *App) initDI(ctx context.Context) {
	a.diContainer = &diContainer{}
}

func (a *App) runConsumer(ctx context.Context) error {
	startService := a.diContainer.AssemblyService(ctx)
	return startService.Process(ctx)
}
