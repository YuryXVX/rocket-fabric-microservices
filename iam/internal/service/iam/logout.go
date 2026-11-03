package iam_service

import (
	"context"
	"iam/internal/service/input"
)

func (s *service) Logout(ctx context.Context, in input.LogoutInput) error {
	if err := s.sessionRepository.Logout(ctx, in); err != nil {
		return err
	}

	return nil
}
