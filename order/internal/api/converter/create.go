package converter

import (
	"order/internal/service/input"
	orderv1 "shared/pkg/openapi/order/v1"

	"github.com/google/uuid"
)

func toPrt(val orderv1.OptNilUUID) *uuid.UUID {
	if !val.Set || val.Null {
		return nil
	}

	u := val.Value
	return &u
}

func CreteRequestToInput(req *orderv1.CreateOrderRequest) *input.CreateOrderInput {
	return &input.CreateOrderInput{
		HullUUID:   req.HullUUID,
		EngineUUID: req.EngineUUID,
		ShieldUUID: toPrt(req.ShieldUUID),
		WeaponUUID: toPrt(req.WeaponUUID),
	}
}
