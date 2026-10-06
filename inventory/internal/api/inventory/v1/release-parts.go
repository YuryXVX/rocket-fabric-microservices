package v1

import (
	"context"
	"inventory/internal/api/converter"
	v1 "shared/pkg/proto/inventory/v1"
)

func (a *api) ReleaseParts(ctx context.Context, req *v1.ReleasePartsRequest) (*v1.ReleasePartsResponse, error) {
	err := a.applicationService.Release(ctx, converter.RequestUUIDsToPartFilter(req.Uuids))

	if err != nil {
		return nil, err
	}

	return &v1.ReleasePartsResponse{}, nil
}
