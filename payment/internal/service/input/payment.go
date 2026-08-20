package input

import (
	"github.com/google/uuid"
	"payment/internal/model"
)

type PayOrderInput struct {
	OrderUUID     uuid.UUID
	PaymentMethod model.PaymentMethod
}
