package iam_service

import (
	"context"
	"errors"
	errs "iam/internal/errors"
	"iam/internal/model"
	"iam/internal/service/input"
	"time"

	"github.com/google/uuid"
)

func (s *service) Login(ctx context.Context, in input.LoginInput) (uuid.UUID, error) {
	user, err := s.userRepository.GetUserByLogin(ctx, in.Login)

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return uuid.Nil, errs.ErrInvalidCredentials
		}
		return uuid.Nil, err
	}

	if err := s.hasher.Compare(in.Password, user.PasswordHash); err != nil {
		return uuid.Nil, errs.ErrInvalidCredentials
	}

	sessionUUID, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, err
	}

	now := time.Now()

	sess := model.Session{
		UUID:      sessionUUID,
		UserUUID:  user.UUID,
		Login:     user.Login,
		CreatedAt: now,
		ExpiresAt: now.Add(24 * time.Hour),
	}

	err = s.sessionRepository.Login(ctx, sess)

	if err != nil {
		return uuid.Nil, err
	}

	return sessionUUID, nil
}
