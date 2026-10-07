package order

type service struct {
	inventoryClient InventoryClientGrpc
	paymentClient   PaymentClientGrpc
	orderRepository OrderRepository
	txManager       TxManager
}

func New(
	inventory InventoryClientGrpc,
	payment PaymentClientGrpc,
	orderRepository OrderRepository,
	txManager TxManager,
) *service {
	return &service{
		inventoryClient: inventory,
		paymentClient:   payment,
		orderRepository: orderRepository,
		txManager:       txManager,
	}
}
