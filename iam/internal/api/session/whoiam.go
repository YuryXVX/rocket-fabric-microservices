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

func (a *api) Whoami(ctx context.Context, req *v1.WhoamiRequest) (*v1.WhoamiResponse, error) {
	input, err := converter.RequestWhoIamToModel(req)

	if err != nil {
		if errors.Is(err, errs.ErrInvalidUUID) {
			return nil, status.Errorf(codes.InvalidArgument, "неверный session uuid")
		}

		return nil, err
	}

	session, err := a.service.WhoIam(ctx, input)

	if err != nil {
		return nil, err
	}

	return converter.SessionModelToResponse(session), nil
}
