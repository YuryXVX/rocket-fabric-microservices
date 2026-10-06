package converter

import (
	"fmt"
	"inventory/internal/model"
	"inventory/internal/model/entity"
	"inventory/internal/model/valueobject"
	"inventory/internal/service/input"
	v1 "shared/pkg/proto/inventory/v1"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToPartTypeModel(t v1.PartType) model.PartType {
	switch t {
	case v1.PartType_PART_TYPE_HULL:
		return model.PartTypeHull
	case v1.PartType_PART_TYPE_ENGINE:
		return model.PartTypeEngine
	case v1.PartType_PART_TYPE_WEAPON:
		return model.PartTypeWeapon
	case v1.PartType_PART_TYPE_SHIELD:
		return model.PartTypeShield
	default:
		return model.PartTypeUnspecified
	}
}

func RequestToInputPartFilter(req *v1.ListPartsRequest) *input.PartFilter {
	uuids := make([]uuid.UUID, len(req.Uuids))

	if len(req.Uuids) > 0 {
		for i, f := range req.Uuids {
			uuids[i] = uuid.MustParse(f)
		}
	}

	return &input.PartFilter{
		UUIDs:    uuids,
		PartType: ToPartTypeModel(req.PartType),
	}
}

func convertModelPartTypeToProto(t valueobject.PartType) v1.PartType {
	switch t {
	case valueobject.PartTypeShield:
		return v1.PartType_PART_TYPE_SHIELD
	case valueobject.PartTypeEngine:
		return v1.PartType_PART_TYPE_ENGINE
	case valueobject.PartTypeHull:
		return v1.PartType_PART_TYPE_HULL
	case valueobject.PartTypeWeapon:
		return v1.PartType_PART_TYPE_WEAPON
	default:
		return v1.PartType_PART_TYPE_UNSPECIFIED
	}
}

func ModelPartToProtoPart(part *entity.Part) *v1.Part {
	return &v1.Part{
		Uuid:          part.UUID().String(),
		Name:          part.Name(),
		Description:   part.Description(),
		Price:         part.Price(),
		StockQuantity: part.StockQuantity(),
		PartType:      convertModelPartTypeToProto(*part.PartType()),
		CreatedAt:     timestamppb.New(part.CreatedAt()),
	}
}

func ModelPartListToProtoPartList(parts []*entity.Part) []*v1.Part {
	protoParts := make([]*v1.Part, len(parts))

	for i, part := range parts {
		protoParts[i] = ModelPartToProtoPart(part)
	}

	return protoParts
}

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func RequestSlotsToFilter(slots *v1.ValidateCompatibilityRequest) (input.PartFilter, error) {
	uuids := make([]uuid.UUID, 0, 4)

	parseAndAppend := func(fieldValue string, fieldName string) error {
		if fieldValue == "" {
			return nil
		}

		parsedUUID, err := parseUUID(fieldValue)

		if err != nil {
			return fmt.Errorf("поле %v имеет валидный uuid %v", fieldName, fieldValue)
		}

		uuids = append(uuids, parsedUUID)

		return nil
	}

	if err := parseAndAppend(slots.EngineUuid, "EngineUUID"); err != nil {
		return input.PartFilter{}, err
	}

	if err := parseAndAppend(slots.HullUuid, "HullUuid"); err != nil {
		return input.PartFilter{}, err
	}

	if err := parseAndAppend(slots.ShieldUuid, "ShieldUuid"); err != nil {
		return input.PartFilter{}, err
	}

	if err := parseAndAppend(slots.WeaponUuid, "WeaponUuid"); err != nil {
		return input.PartFilter{}, err
	}

	return input.PartFilter{
		UUIDs: uuids,
	}, nil
}
