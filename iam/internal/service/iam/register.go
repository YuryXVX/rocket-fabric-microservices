package iam_service

import (
	"context"
	"errors"
	"fmt"
	errs "iam/internal/errors"
	"iam/internal/model"
	user_repository "iam/internal/repository/user"
	"iam/internal/service/input"
	"time"

	"github.com/google/uuid"
)

func (s *service) Register(ctx context.Context, in input.RegisterInput) (uuid.UUID, error) {
	hash, err := s.hasher.Hash(in.Password)

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("ошибка при регистрации пользователя")
	}

	userUUID, err := s.userRepository.Register(ctx, model.User{
		UUID:         uuid.New(),
		Login:        in.Login,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
		UpdatedAt:    nil,
	})

	if err != nil {
		if errors.Is(err, user_repository.ErrUserAlreadyExists) {
			return uuid.UUID{}, errs.ErrUserAlreadyExists
		}

		return uuid.UUID{}, err
	}

	return userUUID, nil
}
