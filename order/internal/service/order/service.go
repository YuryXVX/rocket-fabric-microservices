package order

type service struct {
	inventoryClient InventoryClientGrpc
	paymentClient   PaymentClientGrpc
	orderRepository OrderRepository
	partRepository  PartRepository
	txManager       TxManager
}

func New(
	inventory InventoryClientGrpc,
	payment PaymentClientGrpc,
	orderRepository OrderRepository,
	partRepository PartRepository,
	txManager TxManager,
) *service {
	return &service{
		inventoryClient: inventory,
		paymentClient:   payment,
		orderRepository: orderRepository,
		partRepository:  partRepository,
		txManager:       txManager,
	}
}
