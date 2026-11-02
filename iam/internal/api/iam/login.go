package api_session

import (
	"context"
	v1 "shared/pkg/proto/auth/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method Login not implemented")
}
