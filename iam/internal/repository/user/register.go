package user_repository

import (
	"context"
	"fmt"
	"iam/internal/model"
	"iam/internal/repository/converter"
)

func (r *repository) Register(ctx context.Context, model model.User) {
	fmt.Printf("user %+v", converter.UserModelToRecord(model))
}
