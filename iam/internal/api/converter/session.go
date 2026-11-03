package converter

import (
	errs "iam/internal/errors"
	"iam/internal/model"
	"iam/internal/service/input"
	v1 "shared/pkg/proto/auth/v1"
	v1Common "shared/pkg/proto/common/v1"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Field struct {
	Name  string
	Value string
}

func (f *Field) IsEmptyValue() bool {
	return f.Value == ""
}

func RequestLoginToInput(req *v1.LoginRequest) (input.LoginInput, error) {
	var fields = []Field{
		{Name: "Login", Value: req.Login},
		{Name: "Password", Value: req.Password},
	}

	if req.Login == "" && req.Password == "" {
		return input.LoginInput{}, errs.ErrEmptyCredential
	}

	for _, field := range fields {
		if field.Value == "" {
			if field.Name == "Login" {
				return input.LoginInput{}, errs.ErrEmptyLogin
			}
			if field.Name == "Password" {
				return input.LoginInput{}, errs.ErrEmptyPassword
			}
		}
	}

	return input.LoginInput{
		Login:    req.Login,
		Password: req.Password,
	}, nil
}

func SessionModelToResponse(session model.Session) *v1.WhoamiResponse {
	createdAt := timestamppb.New(session.CreatedAt)
	expiresAt := timestamppb.New(session.ExpiresAt)

	var updateAt timestamppb.Timestamp

	if session.UpdatedAt != nil {
		updateAt = *timestamppb.New(*session.UpdatedAt)
	}

	return &v1.WhoamiResponse{
		Session: &v1Common.Session{
			Uuid:      session.UUID.String(),
			CreatedAt: createdAt,
			UpdatedAt: &updateAt,
			ExpiresAt: expiresAt,
		},
	}
}

func RequestWhoIamToModel(req *v1.WhoamiRequest) (input.WhoIamInput, error) {
	sessionUUID, err := uuid.Parse(req.SessionUuid)

	if err != nil {
		return input.WhoIamInput{}, errs.ErrInvalidUUID
	}

	return input.WhoIamInput{
		SessionUuid: sessionUUID,
	}, nil
}

func RequestLogoutToInput(req *v1.LogoutRequest) (input.LogoutInput, error) {
	sessionUUID, err := uuid.Parse(req.SessionUuid)

	if err != nil {
		return input.LogoutInput{}, errs.ErrInvalidUUID
	}

	return input.LogoutInput{
		SessionUUID: sessionUUID,
	}, nil
}
