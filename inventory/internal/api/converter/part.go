package converter

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"inventory/internal/model"
	"inventory/internal/service/input"
	v1 "shared/pkg/proto/inventory/v1"
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

func convertModelPartTypeToProto(t model.PartType) v1.PartType {
	switch t {
	case model.PartTypeUnspecified:
		return v1.PartType_PART_TYPE_UNSPECIFIED
	case model.PartTypeEngine:
		return v1.PartType_PART_TYPE_ENGINE
	case model.PartTypeHull:
		return v1.PartType_PART_TYPE_HULL
	default:
		return v1.PartType_PART_TYPE_UNSPECIFIED
	}
}

func ModelPartToProtoPart(part *model.Part) *v1.Part {
	return &v1.Part{
		Uuid:          part.UUID.String(),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		PartType:      convertModelPartTypeToProto(part.PartType),
		CreatedAt:     timestamppb.New(part.CreatedAt),
	}
}

func ModelPartListToProtoPartList(parts []*model.Part) []*v1.Part {
	protoParts := make([]*v1.Part, len(parts))

	for i, part := range parts {
		protoParts[i] = ModelPartToProtoPart(part)
	}

	return protoParts
}
