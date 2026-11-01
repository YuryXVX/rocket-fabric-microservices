package api_user

import (
	"context"
	v1 "shared/pkg/proto/user/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) GetUser(context.Context, *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetUser not implemented")
}
