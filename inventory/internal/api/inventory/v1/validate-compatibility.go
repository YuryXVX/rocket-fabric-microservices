package v1

import (
	"context"
	"inventory/internal/api/converter"
	v1 "shared/pkg/proto/inventory/v1"
)

func (a *api) ValidateCompatibility(ctx context.Context, req *v1.ValidateCompatibilityRequest) (*v1.ValidateCompatibilityResponse, error) {
	input, err := converter.RequestSlotsToFilter(req)

	if err != nil {
		return nil, err
	}

	err = a.applicationService.ValidateCompatibility(ctx, input)

	if err != nil {
		return nil, err
	}

	return &v1.ValidateCompatibilityResponse{}, nil
}
