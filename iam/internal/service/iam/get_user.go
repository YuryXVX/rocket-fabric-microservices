package iam_service

import (
	"context"
	"errors"
	errs "iam/internal/errors"
	"iam/internal/model"
	user_repository "iam/internal/repository/user"

	"github.com/google/uuid"
)

func (s *service) GetUser(ctx context.Context, userUUID uuid.UUID) (model.User, error) {
	user, err := s.userRepository.GetUser(ctx, userUUID)

	if err != nil {
		if errors.Is(err, user_repository.ErrUserNotFound) {
			return model.User{}, errs.ErrUserNotFound
		}
		return model.User{}, err
	}

	return user, nil
}
