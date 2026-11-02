package api_user

import (
	"context"
	"iam/internal/api/converter"
	v1 "shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	uuid, err := a.service.Register(ctx, converter.RequestToInput(req))

	if err != nil {
		return nil, err
	}

	return &v1.RegisterResponse{
		UserUuid: uuid.String(),
	}, nil
}
