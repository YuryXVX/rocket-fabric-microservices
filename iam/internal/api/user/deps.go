package api_user

import (
	"context"
	"iam/internal/model"
	"iam/internal/service/input"

	"github.com/google/uuid"
)

type UserService interface {
	Register(context.Context, input.RegisterInput) (uuid.UUID, error)
	GetUser(context.Context, uuid.UUID) (model.User, error)
}
