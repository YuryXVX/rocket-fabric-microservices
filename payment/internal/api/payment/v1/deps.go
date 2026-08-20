package v1

import (
	"context"

	"github.com/google/uuid"
	"payment/internal/service/input"
)

type PayService interface {
	Pay(ctx context.Context, in input.PayOrderInput) (uuid.UUID, error)
}
