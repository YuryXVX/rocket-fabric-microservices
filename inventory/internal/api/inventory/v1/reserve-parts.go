package v1

import (
	"context"
	"inventory/internal/api/converter"
	v1 "shared/pkg/proto/inventory/v1"
)

func (a *api) ReserveParts(ctx context.Context, req *v1.ReservePartsRequest) (*v1.ReservePartsResponse, error) {
	err := a.applicationService.Reserve(ctx, converter.RequestUUIDsToPartFilter(req.Uuids))

	if err != nil {
		return nil, err
	}

	return nil, nil
}
