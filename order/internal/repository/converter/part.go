package converter

import (
	"order/internal/model"
	"order/internal/repository/record"
)

func OrderItemToRecord(item model.OrderItem) record.PartRecord {
	return record.PartRecord{
		PartUUID: item.PartUUID,
		Price:    item.Price,
		PartType: record.PartType(item.PartType),
	}
}

func RecordItemToModel(item record.PartRecord) model.OrderItem {
	return model.OrderItem{
		PartUUID: item.PartUUID,
		Price:    item.Price,
		PartType: model.PartType(item.PartType),
	}
}

func OrderItemsToRecords(orderItems []model.OrderItem) []record.PartRecord {
	r := make([]record.PartRecord, 0, len(orderItems))

	for _, item := range orderItems {
		r = append(r, OrderItemToRecord(item))
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
