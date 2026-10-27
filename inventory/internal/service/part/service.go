package inventory

type service struct {
	inventoryRepository InventoryRepository
}

func New(inventoryRepository InventoryRepository) *service {
	return &service{inventoryRepository: inventoryRepository}
}
