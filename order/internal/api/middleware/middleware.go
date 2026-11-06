package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"order/internal/model"
	"plaform/pkg/auth"
	"time"

	"github.com/google/uuid"
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
	WhoIam(context.Context, string) (model.Session, error)
}

func AuthMiddleware(iam IamClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")

			if h == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			session, err := iam.WhoIam(r.Context(), h)

			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			uuidUser, err := uuid.Parse(session.UserUUID)

			if err != nil {
				http.Error(w, "Ошибка формата uuid пользователя", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = auth.WithUserUUID(ctx, uuidUser)
			ctx = auth.WithSessionUUID(ctx, session.SessionUUID)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}
