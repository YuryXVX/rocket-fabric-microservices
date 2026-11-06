package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"plaform/pkg/auth"
	"time"
)

// RequestLogger создает middleware для логирования времени выполнения запросов
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Засекаем время начала обработки запроса
		startTime := time.Now()

		// Логируем начало запроса
		slog.InfoContext(r.Context(), "⏱️ начало запроса", "method", r.Method, "path", r.URL.Path)

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)

		// Вычисляем время выполнения запроса
		duration := time.Since(startTime)

		// Логируем окончание запроса с временем выполнения
		slog.InfoContext(r.Context(), "✅ запрос завершен", "method", r.Method, "path", r.URL.Path, "duration", duration)
	})
}

type IamClient interface {
	WhoIam(context.Context, string) (string, error)
}

func AuthMiddleware(iam IamClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")

			if h == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			sessionUUID, err := iam.WhoIam(r.Context(), h)

			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(
				w,
				r.WithContext(auth.WithSessionUUID(r.Context(), sessionUUID)),
			)
		})
	}
}
