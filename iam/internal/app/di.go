package app

import (
	"context"
	apiAuthV1 "iam/internal/api/iam"
	api_user "iam/internal/api/user"
	"plaform/pkg/di"
	v1Auth "shared/pkg/proto/auth/v1"
	v1User "shared/pkg/proto/user/v1"
)

type diContainer struct {
	authHandlers di.Value[v1Auth.AuthServiceServer]
	userHandlers di.Value[v1User.UserServiceServer]
}

func (di *diContainer) AuthHandlers(ctx context.Context) v1Auth.AuthServiceServer {
	return di.authHandlers.Get(ctx, func(ctx context.Context) v1Auth.AuthServiceServer {
		return apiAuthV1.New()
	})
}

func (di *diContainer) UserHandlers(ctx context.Context) v1User.UserServiceServer {
	return di.userHandlers.Get(ctx, func(ctx context.Context) v1User.UserServiceServer {
		return api_user.New()
	})
}
