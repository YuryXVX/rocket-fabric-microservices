package inventory

import (
	"context"
	"order/internal/client/grpc/inventory/converter"
	"order/internal/model"
	v1 "shared/pkg/proto/inventory/v1"

	"github.com/google/uuid"
)

type inventoryClient struct {
	grpc v1.PartServiceClient
}

func New(c v1.PartServiceClient) *inventoryClient {
	return &inventoryClient{
		grpc: c,
	}
}

func (i *inventoryClient) ListParts(ctx context.Context, in []uuid.UUID) ([]model.OrderItem, error) {
	parts, err := i.grpc.ListParts(ctx, converter.InputToListRequest(in))

	if err != nil {
		return []model.OrderItem{}, err
	}

	return converter.ListResponseToOrderItem(parts), nil
}
