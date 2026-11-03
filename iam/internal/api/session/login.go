package api_session

import (
	"context"
	v1 "shared/pkg/proto/auth/v1"

	"github.com/google/uuid"
)

func (a *api) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	return &v1.LoginResponse{
		SessionUuid: uuid.UUID{}.String(),
	}, nil
}
