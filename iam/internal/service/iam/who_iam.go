package iam_service

import (
	"context"
	"errors"
	errs "iam/internal/errors"
	"iam/internal/model"
	sessionRepository "iam/internal/repository/session"
	"iam/internal/service/input"
)

func (s *service) WhoIam(ctx context.Context, in input.WhoIamInput) (model.Session, error) {
	session, err := s.sessionRepository.WhoIam(ctx, in.SessionUuid)

	if err != nil {
		if errors.Is(err, sessionRepository.ErrSessionNotFound) {
			return model.Session{}, errs.ErrSessionNotFound
		}

		return model.Session{}, err
	}

	return session, nil
}
