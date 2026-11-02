package api_session

import (
	"context"
	"iam/internal/model"
	"iam/internal/service/input"

	"github.com/google/uuid"
)

type SessionService interface {
	Login(context.Context, input.LoginInput) (uuid.UUID, error)
	Logout(context.Context, input.LogoutInput) error
	WhoIam(context.Context, input.WhoIamInput) (model.Session, error)
}
