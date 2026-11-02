package iam_service

import (
	"context"
	"iam/internal/model"
	"iam/internal/service/input"
)

func (s *service) WhoIam(ctx context.Context, in input.WhoIamInput) (model.Session, error) {
	return model.Session{}, nil
}
