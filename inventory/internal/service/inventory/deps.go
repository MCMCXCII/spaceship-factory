package inventory

import (
	"context"

	"github.com/MCMCXCII/spaceship-factory/inventory/internal/model"
	"github.com/google/uuid"
)

type InventoryRepository interface {
	Get(context.Context, uuid.UUID) (model.Part, error)
	List(context.Context) ([]model.Part, error)
}
