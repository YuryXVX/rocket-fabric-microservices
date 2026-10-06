package application

type service struct {
	inventoryRepository  InventoryRepository
	compatibilityChecker CompatibilityChecker
}

func NewService(
	inventoryRepository InventoryRepository,
	compatibilityChecker CompatibilityChecker,
) *service {
	return &service{
		inventoryRepository:  inventoryRepository,
		compatibilityChecker: compatibilityChecker,
	}
}
