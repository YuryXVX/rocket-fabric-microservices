package application

type service struct {
	inventoryRepository  InventoryRepository
	compatibilityChecker CompatibilityChecker
	txManager            TxManager
}

func NewService(
	inventoryRepository InventoryRepository,
	compatibilityChecker CompatibilityChecker,
	txManager TxManager,
) *service {
	return &service{
		inventoryRepository:  inventoryRepository,
		compatibilityChecker: compatibilityChecker,
		txManager:            txManager,
	}
}
