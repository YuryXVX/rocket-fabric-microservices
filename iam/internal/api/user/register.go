package api_user

import (
	"context"
	"errors"
	"iam/internal/api/converter"
	errs "iam/internal/errors"
	v1 "shared/pkg/proto/user/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	payload, err := converter.RequestRegisterToInput(req)

	if err != nil {
		if errors.Is(err, errs.ErrWeakPassword) || errors.Is(err, errs.ErrEmptyCredential) {
			return nil, status.Errorf(codes.InvalidArgument, "неверный логин или пароль")
		}

		return nil, err
	}

	userUUID, err := a.service.Register(ctx, payload)

	if err != nil {
		if errors.Is(err, errs.ErrUserAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "пользователь с  логином %v уже сущесвутет", req.GetInfo().Info.GetLogin())
		}

		return nil, err
	}

	return &v1.RegisterResponse{
		UserUuid: userUUID.String(),
	}, nil
}
