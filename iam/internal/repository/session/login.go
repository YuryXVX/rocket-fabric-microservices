package session_repository

import (
	"context"
	"fmt"
	"iam/internal/model"
	"time"
)

func (r *repository) Login(ctx context.Context, sess model.Session) error {
	key := fmt.Sprintf("session:%s", sess.UUID.String())

	// Вычисляем TTL на основе переданного времени истечения
	ttl := time.Until(sess.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("session has already expired")
	}

	pipe := r.client.TxPipeline()

	// Записываем типизированные поля в HSET (go-redis автоматически преобразует базовые типы в строки)
	pipe.HSet(ctx, key, map[string]interface{}{
		"user_uuid":  sess.UserUUID.String(),
		"login":      sess.Login,
		"created_at": sess.CreatedAt.Format(time.RFC3339),
		"expires_at": sess.ExpiresAt.Format(time.RFC3339),
	})

	// Выставляем TTL для ключа в Redis
	pipe.Expire(ctx, key, ttl)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute redis pipeline for session: %w", err)
	}

	return nil
}
