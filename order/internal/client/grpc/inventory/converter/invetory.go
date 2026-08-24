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
			PartType: model.PartType(part.PartType.String()),
			Price:    part.Price,
		})
	}

	return items
}
