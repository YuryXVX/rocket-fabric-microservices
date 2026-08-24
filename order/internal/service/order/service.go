package order

type service struct {
	inventoryClient InventoryClientGrpc
	paymentClient   PaymentClientGrpc
}

func New(
	inventory InventoryClientGrpc,
	payment PaymentClientGrpc,
) *service {
	return &service{
		inventoryClient: inventory,
		paymentClient:   payment,
	}
}
