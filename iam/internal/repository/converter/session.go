package converter

import (
	"iam/internal/model"
	"iam/internal/repository/redis_view"
)

func SessionModelToRedisView(session model.Session) redis_view.SessionView {
	return redis_view.SessionView{
		UUID:      session.UUID,
		UserUUID:  session.UserUUID,
		Login:     session.Login,
		CreatedAt: session.CreatedAt,
		ExpiresAt: session.ExpiresAt,
	}
}

func SessionRedisViewToModel(view redis_view.SessionView) model.Session {
	return model.Session{
		UUID:      view.UUID,
		UserUUID:  view.UserUUID,
		Login:     view.Login,
		CreatedAt: view.CreatedAt,
		ExpiresAt: view.ExpiresAt,
	}
}
