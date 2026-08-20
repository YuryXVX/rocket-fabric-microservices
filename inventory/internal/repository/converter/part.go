package converter

import (
	"github.com/google/uuid"
	"inventory/internal/model"
	"inventory/internal/repository/record"
)

func ConvertPartType(t record.PartType) model.PartType {
	switch t {
	case record.PartTypeUnspecified:
		return model.PartTypeUnspecified
	case record.PartTypeHull:
		return model.PartTypeHull
	case record.PartTypeEngine:
		return model.PartTypeEngine
	case record.PartTypeShield:
		return model.PartTypeShield
	case record.PartTypeWeapon:
		return model.PartTypeWeapon
	default:
		return model.PartTypeUnspecified
	}
}

func ConvertRecordPartType(t model.PartType) record.PartType {
	switch t {
	case model.PartTypeHull:
		return record.PartTypeHull
	case model.PartTypeEngine:
		return record.PartTypeEngine
	case model.PartTypeShield:
		return record.PartTypeShield
	case model.PartTypeWeapon:
		return record.PartTypeWeapon
	default:
		return record.PartTypeUnspecified
	}
}

func RecordPartToModel(record record.Part) *model.Part {
	return &model.Part{
		UUID:          uuid.MustParse(record.UUID),
		Name:          record.Name,
		Description:   record.Description,
		Price:         record.Price,
		PartType:      ConvertPartType(record.PartType),
		StockQuantity: int64(record.StockQuantity),
		CreatedAt:     record.CreatedAt,
	}
}
