package main

import (
	"context"
	"inventory/internal/app"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func MustLoadConfig() {
	err := godotenv.Load("../inventory.env")

	if err != nil {
		slog.Error("ошибка загрузки переменных окружения из inventory.env", "error", err)
		os.Exit(1)
	}
}

func main() {
	MustLoadConfig()

	a := app.New(context.Background())

	if err := a.Run(); err != nil {
		slog.Error("ошибка при работе приложения", "error", err)
	}
}
