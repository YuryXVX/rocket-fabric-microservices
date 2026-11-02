package api_session

import (
	"context"
	v1 "shared/pkg/proto/auth/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) Logout(ctx context.Context, req *v1.LogoutRequest) (*v1.LogoutResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method Logout not implemented")
}
