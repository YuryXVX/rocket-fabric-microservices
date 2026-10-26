package main

import (
	"assemply/internal/app"
	"context"
	"log/slog"
	"os"
)

func main() {
	app := app.NewApp(context.Background())

	if err := app.Run(); err != nil {
		slog.Error("не удалось запустить приложение", "error", err)

		os.Exit(1)
		return
	}
}
