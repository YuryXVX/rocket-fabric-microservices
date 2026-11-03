package api_session

import (
	"context"
	v1 "shared/pkg/proto/auth/v1"
	v1Common "shared/pkg/proto/common/v1"
)

func (a *api) Whoami(ctx context.Context, req *v1.WhoamiRequest) (*v1.WhoamiResponse, error) {
	return &v1.WhoamiResponse{
		Session: &v1Common.Session{},
	}, nil
}
