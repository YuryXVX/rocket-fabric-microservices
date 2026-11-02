package redis_view

import (
	"time"

	"github.com/google/uuid"
)

type SessionView struct {
	UUID      uuid.UUID `redis:"uuid"`
	UserUUID  uuid.UUID `redis:"user_uuid"`
	Login     string    `redis:"login"`
	CreatedAt time.Time `redis:"created_at"`
	ExpiresAt time.Time `redis:"expires_at"`
}
