package api_user

import (
	"context"
	v1 "shared/pkg/proto/user/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) Register(context.Context, *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method Register not implemented")
}
