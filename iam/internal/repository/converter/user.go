package converter

import (
	"iam/internal/model"
	records "iam/internal/repository/record"
)

func UserModelToRecord(model model.User) records.UserRecord {
	return records.UserRecord{
		UUID:         model.UUID,
		Login:        model.Login,
		PasswordHash: model.PasswordHash,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}

func RecordToUserModel(rec records.UserRecord) model.User {
	return model.User{
		UUID:         rec.UUID,
		Login:        rec.Login,
		PasswordHash: rec.PasswordHash,
		CreatedAt:    rec.CreatedAt,
		UpdatedAt:    rec.UpdatedAt,
	}
}
