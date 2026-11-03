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

func (a *api) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	payload, err := converter.RequestLoginToInput(req)

	if err != nil {
		if errors.Is(err, errs.ErrEmptyLogin) {
			return nil, status.Errorf(codes.InvalidArgument, "логин должен быть заполнен")
		}

		if errors.Is(err, errs.ErrEmptyPassword) {
			return nil, status.Errorf(codes.InvalidArgument, "пароль должен быть заполнен")
		}

		return nil, err
	}

	sessionUUID, err := a.service.Login(ctx, payload)

	if err != nil {
		return nil, err
	}

	return &v1.LoginResponse{
		SessionUuid: sessionUUID.String(),
	}, nil
}
