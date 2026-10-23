package v1

import (
	"context"

	"github.com/MCMCXCII/spaceship-factory/inventory/internal/model"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/service/input"
	"github.com/google/uuid"
)

type InventoryService interface {
	Get(context.Context, uuid.UUID) (model.Part, error)
	List(context.Context, input.PartFilter) ([]model.Part, error)
}
