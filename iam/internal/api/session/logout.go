package api_session

import (
	"context"
	v1 "shared/pkg/proto/auth/v1"
)

func (a *api) Logout(ctx context.Context, req *v1.LogoutRequest) (*v1.LogoutResponse, error) {
	return &v1.LogoutResponse{}, nil
}
