package iam_service

type service struct {
	userRepository    UserRepository
	sessionRepository SessionRepository
	hasher            Hasher
}

func New(
	userRepo UserRepository,
	sessionRepo SessionRepository,
	hasher Hasher,
) *service {
	return &service{
		userRepository:    userRepo,
		sessionRepository: sessionRepo,
		hasher:            hasher,
	}
}
