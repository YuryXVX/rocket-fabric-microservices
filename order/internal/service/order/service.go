package order

type service struct {
	inventoryClient InventoryClientGrpc
	paymentClient   PaymentClientGrpc
	orderRepository OrderRepository
}

func New(
	inventory InventoryClientGrpc,
	payment PaymentClientGrpc,
	orderRepository OrderRepository,
) *service {
	return &service{
		inventoryClient: inventory,
		paymentClient:   payment,
		orderRepository: orderRepository,
	}
}
