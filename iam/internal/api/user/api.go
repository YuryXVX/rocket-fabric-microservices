package api_user

import (
	v1 "shared/pkg/proto/user/v1"
)

type api struct {
	v1.UnimplementedUserServiceServer
}

func New() *api {
	return &api{}
}
