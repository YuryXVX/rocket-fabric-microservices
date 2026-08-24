package order

import (
	"context"
	"order/internal/model"

	"github.com/google/uuid"
)

func (s *service) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	return model.Order{}, nil
}
