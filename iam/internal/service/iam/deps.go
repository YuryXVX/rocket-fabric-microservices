package iam_service

import (
	"context"
	"iam/internal/model"
	"iam/internal/service/input"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Login(context.Context, input.LoginInput) (uuid.UUID, error)
	Logout(context.Context, input.LogoutInput) error
	WhoIam(context.Context, uuid.UUID) (model.Session, error)
}

type UserRepository interface {
	Register(context.Context, model.User) (uuid.UUID, error)
	GetUser(context.Context, uuid.UUID) (model.User, error)
}

type Hasher interface {
	Hash(string) (string, error)
	Compare(password, hash string) error
}
