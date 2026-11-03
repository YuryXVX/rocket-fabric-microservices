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

func (a *api) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	userUUID, err := converter.RequestGetUserToInput(req)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "неверный формат uuid пользователя")
	}

	user, err := a.service.GetUser(ctx, userUUID)

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "пользователь с UUID %v не найден", req.UserUuid)
		}
		return nil, err
	}

	return converter.UserModelToRequest(user), nil
}
