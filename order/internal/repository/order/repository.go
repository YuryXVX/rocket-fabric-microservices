package order

type repository struct{}

func New() *repository {
	return &repository{}
}
