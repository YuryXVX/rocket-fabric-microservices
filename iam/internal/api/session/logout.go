package api_session

import (
	"context"
	"errors"
	"iam/internal/api/converter"
	errs "iam/internal/errors"
	v1 "shared/pkg/proto/auth/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) Logout(ctx context.Context, req *v1.LogoutRequest) (*v1.LogoutResponse, error) {
	input, err := converter.RequestLogoutToInput(req)

	if err != nil {
		if errors.Is(err, errs.ErrInvalidUUID) {
			return nil, status.Errorf(codes.InvalidArgument, "не валидный session uuid")
		}
		return nil, err
	}

	if err := a.service.Logout(ctx, input); err != nil {
		return nil, err
	}

	return &v1.LogoutResponse{}, nil
}
