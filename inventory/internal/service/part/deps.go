package inventory

import (
	"context"

	"github.com/google/uuid"

	"github.com/MCMCXCII/spaceship-factory/inventory/internal/model"
)

type InventoryRepository interface {
	Get(context.Context, uuid.UUID) (model.Part, error)
	List(context.Context) ([]model.Part, error)
}
