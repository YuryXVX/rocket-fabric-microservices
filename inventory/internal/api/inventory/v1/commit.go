package v1

import (
	"context"
	"errors"
	"inventory/internal/api/converter"
	errs "inventory/internal/errors"
	v1 "shared/pkg/proto/inventory/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) CommitParts(ctx context.Context, req *v1.CommitPartsRequest) (*v1.CommitPartsResponse, error) {
	err := a.applicationService.CommitParts(ctx, converter.RequestToInputCommitParts(req))

	if err != nil {
		if errors.Is(err, errs.ErrPartNotFound) {
			return nil, status.Error(codes.NotFound, "one or more requested parts were not found in the inventory")
		}
		if errors.Is(err, errs.ErrPartsNotAvailable) {
			return nil, status.Error(codes.FailedPrecondition, "insufficient stock_quantity or reserved amount for some parts")
		}

		return nil, status.Errorf(codes.Internal, "internal database error: %v", err)

	}

	return &v1.CommitPartsResponse{}, nil
}
