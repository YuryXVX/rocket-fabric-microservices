package converter

import (
	"order/internal/model"
	"order/internal/repository/record"

	"github.com/google/uuid"
)

func OrderItemTypeToRecord(t model.PartType) record.PartType {
	switch t {
	case model.PartTypeWeapon:
		return record.PartTypeWeapon
	case model.PartTypeEngine:
		return record.PartTypeEngine
	case model.PartTypeShield:
		return record.PartTypeShield
	case model.PartTypeHull:
		return record.PartTypeHull
	default:
		return record.PartTypeUnspecified
	}
}

func OrderItemToRecord(item model.OrderItem, orderUUID uuid.UUID) record.PartRecord {
	return record.PartRecord{
		PartUUID: item.PartUUID,
		Price:    item.Price,
		PartType: OrderItemTypeToRecord(item.PartType),
	}
}

func RecordItemToModel(item record.PartRecord) model.OrderItem {
	return model.OrderItem{
		PartUUID: item.PartUUID,
		Price:    item.Price,
		PartType: model.PartType(item.PartType),
	}
}

func OrderItemsToRecords(orderItems []model.OrderItem, orderUUID uuid.UUID) []record.PartRecord {
	r := make([]record.PartRecord, 0, len(orderItems))

	for _, item := range orderItems {
		r = append(r, OrderItemToRecord(item, orderUUID))
	}

	return r
}

func RecordItemsToModel(orderItems []record.PartRecord) []model.OrderItem {
	r := make([]model.OrderItem, 0, len(orderItems))

	for _, item := range orderItems {
		r = append(r, RecordItemToModel(item))
	}

	return r
}
