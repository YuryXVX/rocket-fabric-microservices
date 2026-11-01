package api_iam

import (
	v1Auth "shared/pkg/proto/auth/v1"
)

type api struct {
	v1Auth.UnimplementedAuthServiceServer
}

func New() *api {
	return &api{}
}
