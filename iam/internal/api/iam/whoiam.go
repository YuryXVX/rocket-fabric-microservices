package api_session

import (
	"context"
	v1 "shared/pkg/proto/auth/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) Whoami(ctx context.Context, req *v1.WhoamiRequest) (*v1.WhoamiResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method Whoami not implemented")
}
