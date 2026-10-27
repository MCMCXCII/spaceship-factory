package order

import (
	"sync"

	"github.com/google/uuid"

	"github.com/MCMCXCII/spaceship-factory/order/internal/repository/record"
)

type Repository struct {
	mu   sync.RWMutex
	data map[uuid.UUID]record.Order
}

func New() *Repository {
	return &Repository{data: make(map[uuid.UUID]record.Order)}
}
