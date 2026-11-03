package converter

import (
	errs "iam/internal/errors"
	"iam/internal/model"
	"iam/internal/service/input"
	v1Common "shared/pkg/proto/common/v1"
	v1 "shared/pkg/proto/user/v1"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func RequestRegisterToInput(req *v1.RegisterRequest) (input.RegisterInput, error) {
	if req.Info == nil {
		return input.RegisterInput{}, errs.ErrInvalidRegister
	}

	if len(req.Info.Password) <= 5 {
		return input.RegisterInput{}, errs.ErrWeakPassword
	}

	if req.Info.Info.Login == "" {
		return input.RegisterInput{}, errs.ErrEmptyCredential
	}

	return input.RegisterInput{
		Login:    req.Info.Info.GetLogin(),
		Password: req.Info.GetPassword(),
	}, nil
}

func RequestGetUserToInput(req *v1.GetUserRequest) (uuid.UUID, error) {
	userUUID, err := uuid.Parse(req.UserUuid)

	if err != nil {
		return uuid.UUID{}, err
	}

	return userUUID, nil
}

func UserModelToRequest(user model.User) *v1.GetUserResponse {
	var updatedAt *timestamppb.Timestamp

	return &v1.GetUserResponse{
		User: &v1Common.User{
			Uuid:      user.UUID.String(),
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: updatedAt,
			Info: &v1Common.UserInfo{
				Login: user.Login,
			},
		},
	}
}
