package api_session

import (
	v1Auth "shared/pkg/proto/auth/v1"
)

type api struct {
	v1Auth.UnimplementedAuthServiceServer

	service SessionService
}

func New(sessionService SessionService) *api {
	return &api{
		service: sessionService,
	}
}
