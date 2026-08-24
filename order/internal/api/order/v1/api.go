package v1

type api struct {
	service OrderService
}

func New(s OrderService) *api {
	return &api{
		service: s,
	}
}
