package main

import (
	"context"
	"log/slog"
	"order/internal/app"

	"github.com/joho/godotenv"
)

func MustLoadEnv() {
	err := godotenv.Load("../order.env")

	if err != nil {
		slog.Error("ошибка загрузки переменных окружения из order.env", "error", err)

		return
	}
}

func main() {
	MustLoadEnv()

	a := app.New(context.Background())

	if err := a.Run(); err != nil {
		slog.Error("ошибка при работе приложения", "error", err)
	}
}
