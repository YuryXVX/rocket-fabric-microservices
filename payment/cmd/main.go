package main

import (
	"context"
	"log/slog"
	"payment/internal/app"
)

func main() {
	a := app.New(context.Background())

	if err := a.Run(); err != nil {
		slog.Error("ошибка при работе приложения", "error", err)
	}
}
