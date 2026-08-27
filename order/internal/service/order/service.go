package order

type service struct {
	inventoryClient InventoryClientGrpc
	paymentClient   PaymentClientGrpc
	orderRepository OrderRepository
	partRepository  PartRepository
}

func New(
	inventory InventoryClientGrpc,
	payment PaymentClientGrpc,
	orderRepository OrderRepository,
	partRepository PartRepository,
) *service {
	return &service{
		inventoryClient: inventory,
		paymentClient:   payment,
		orderRepository: orderRepository,
		partRepository:  partRepository,
	}
}
