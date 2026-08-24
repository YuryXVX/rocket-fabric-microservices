package v1

import (
	"context"
	"errors"

	"inventory/internal/api/converter"
	errs "inventory/internal/errors"
	v1 "shared/pkg/proto/inventory/v1"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) ListParts(ctx context.Context, req *v1.ListPartsRequest) (*v1.ListPartsResponse, error) {
	if len(req.Uuids) > 0 {
		for _, u := range req.Uuids {
			if _, err := uuid.Parse(u); err != nil {
				return &v1.ListPartsResponse{}, status.Errorf(codes.InvalidArgument, "uuid не правильного формата %s", err)
			}
		}
	}

	parts, err := a.serviceInventory.List(ctx, *converter.RequestToInputPartFilter(req))

	if err != nil {
		if errors.Is(err, errs.ErrInvalidUUID) {
			return &v1.ListPartsResponse{}, status.Errorf(codes.InvalidArgument, "uuid не правильного формата %s", err)
		}

		return nil, err
	}

	return &v1.ListPartsResponse{
		Parts: converter.ModelPartListToProtoPartList(parts),
	}, nil
}
