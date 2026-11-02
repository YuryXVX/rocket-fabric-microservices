package api_user

import (
	v1 "shared/pkg/proto/user/v1"
)

type api struct {
	v1.UnimplementedUserServiceServer

	service UserService
}

func New(iam UserService) *api {
	return &api{
		service: iam,
	}
}
