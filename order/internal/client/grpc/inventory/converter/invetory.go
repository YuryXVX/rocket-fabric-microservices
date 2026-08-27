package converter

import (
	"order/internal/model"
	v1 "shared/pkg/proto/inventory/v1"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func InputToListRequest(in []uuid.UUID) *v1.ListPartsRequest {
	uuids := make([]string, 0, len(in))

	for _, uuid := range in {
		uuids = append(uuids, *proto.String(uuid.String()))
	}

	return &v1.ListPartsRequest{
		Uuids: uuids,
	}
}

func ListResponseToOrderItem(res *v1.ListPartsResponse) []model.OrderItem {
	items := make([]model.OrderItem, 0, len(res.Parts))

	for _, part := range res.Parts {
		items = append(items, model.OrderItem{
			PartUUID: uuid.MustParse(part.Uuid),
			PartType: PartTypeOrder(part.PartType),
			Price:    part.Price,
		})
	}

	return items
}

func PartTypeOrder(t v1.PartType) model.PartType {
	var partType model.PartType

	switch t {
	case v1.PartType_PART_TYPE_ENGINE:
		partType = model.PartTypeEngine
	case v1.PartType_PART_TYPE_HULL:
		partType = model.PartTypeHull
	case v1.PartType_PART_TYPE_SHIELD:
		partType = model.PartTypeShield
	case v1.PartType_PART_TYPE_WEAPON:
		partType = model.PartTypeWeapon
	default:
		partType = model.PartTypeUnspecified
	}

	return partType
}
