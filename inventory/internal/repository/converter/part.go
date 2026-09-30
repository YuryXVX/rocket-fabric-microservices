package converter

import (
	"encoding/json"
	"fmt"
	"inventory/internal/model"
	"inventory/internal/model/entity"
	"inventory/internal/model/valueobject"
	record "inventory/internal/repository/record"
)

func RecordPartToModel(record record.Part) *model.Part {
	return &model.Part{
		UUID:          record.UUID,
		Name:          record.Name,
		Description:   record.Description,
		Price:         record.Price,
		PartType:      "",
		StockQuantity: int64(record.StockQuantity),
		CreatedAt:     record.CreatedAt,
	}
}

func RecordPartToDomain(rec record.Part) (entity.Part, error) {
	var partProps record.PartPropertiesRecord

	if err := json.Unmarshal(rec.Properties, &partProps); err != nil {
		return entity.Part{}, fmt.Errorf("десериализовать свойства: %w", err)
	}

	partType, err := valueobject.NewPartType(rec.PartType)

	if err != nil {
		return entity.Part{}, fmt.Errorf("получение типа детали %v", &err)
	}

	props, err := partPropertiesFromRecord(partProps)

	if err != nil {
		return entity.Part{}, fmt.Errorf("десериализация свойства в доменную модель %v", &err)
	}

	return entity.RestorePart(
		rec.UUID,
		rec.Name,
		rec.Description,
		&partType,
		rec.Price,
		rec.StockQuantity,
		rec.Reserved,
		props,
		rec.CreatedAt,
	), nil
}

func partPropertiesFromRecord(rec record.PartPropertiesRecord) (*valueobject.PartProperties, error) {
	switch {
	case rec.Hull != nil:
		return valueobject.NewHullProperties(rec.Hull.Strength)
	case rec.Engine != nil:
		return valueobject.NewEngineProperties(valueobject.EngineClass(rec.Engine.Class), rec.Engine.RequiredStrength)
	case rec.Shield != nil:
		return valueobject.NewShieldProperties(valueobject.ShieldType(rec.Shield.ShieldType))
	case rec.Weapon != nil:
		return valueobject.NewWeaponProperties(valueobject.WeaponType(rec.Weapon.WeaponType))
	default:
		return &valueobject.PartProperties{}, nil
	}
}
