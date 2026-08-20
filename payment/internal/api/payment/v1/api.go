package v1

import v1 "shared/pkg/proto/payment/v1"

type api struct {
	v1.UnimplementedBillingServiceServer

	service PayService
}

func New(s PayService) *api {
	return &api{
		service: s,
	}
}
