package v1

import v1 "shared/pkg/proto/inventory/v1"

type api struct {
	v1.UnimplementedPartServiceServer

	serviceInventory ServiceInventory
}

func New(s ServiceInventory) *api {
	return &api{
		serviceInventory: s,
	}
}
