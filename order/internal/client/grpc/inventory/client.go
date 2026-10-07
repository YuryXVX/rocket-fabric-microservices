package inventory

import (
	"context"
	"order/internal/client/grpc/inventory/converter"
	"order/internal/model"
	"order/internal/service/input"
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

func (i *inventoryClient) ValidateCompatibility(ctx context.Context, in input.CreateOrderInput) error {
	_, err := i.grpc.ValidateCompatibility(ctx, converter.InputPartSlotsToRequest(in))

	return err
}

func (i *inventoryClient) ReserveParts(ctx context.Context, in []uuid.UUID) error {
	_, err := i.grpc.ReserveParts(ctx, converter.InputUUIDToReserveRequest(in))

	return err
}
func (i *inventoryClient) ReleaseParts(ctx context.Context, in []uuid.UUID) error {
	_, err := i.grpc.ReleaseParts(ctx, converter.InputUUIDToReleaseRequest(in))

	return err
}
