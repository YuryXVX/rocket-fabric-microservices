package v1

import v1 "shared/pkg/proto/inventory/v1"

type api struct {
	v1.UnimplementedPartServiceServer

	applicationService ApplicationService
}

func New(s ApplicationService) *api {
	return &api{
		applicationService: s,
	}
}
