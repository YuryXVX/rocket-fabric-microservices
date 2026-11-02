package iam_service

import (
	"context"
	"iam/internal/model"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Login(context.Context, uuid.UUID) (uuid.UUID, error)
	Logout(context.Context, uuid.UUID) error
	WhoIam(context.Context, uuid.UUID) (model.Session, error)
}

type UserRepository interface {
	Register(context.Context, uuid.UUID) (model.User, error)
	GetUser(context.Context, uuid.UUID) (model.User, error)
}
